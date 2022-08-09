package protocol

import (
	"context"
	"net"
	"time"

	"github.com/infraboard/mcenter/apps/instance"
	"github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcenter/client/rpc/lifecycle"
	"google.golang.org/grpc"

	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/grpc/middleware/recovery"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"

	"codeup.aliyun.com/baber/go/keyauth/conf"
)

// NewGRPCService todo
func NewGRPCService() *GRPCService {
	log := zap.L().Named("GRPC Service")

	rc := recovery.NewInterceptor(recovery.NewZapRecoveryHandler())
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		rc.UnaryServerInterceptor(),
	))

	// 控制GRPC启动的其他服务
	ctx, cancel := context.WithCancel(context.Background())

	return &GRPCService{
		svr: grpcServer,
		l:   log,
		c:   conf.C(),

		ctx:    ctx,
		cancel: cancel,
	}
}

// GRPCService grpc服务
type GRPCService struct {
	svr *grpc.Server
	l   logger.Logger
	c   *conf.Config

	ctx    context.Context
	cancel context.CancelFunc
	// 控制实例的上线和下线
	lifecycler lifecycle.Lifecycler
}

// Start 启动GRPC服务
func (s *GRPCService) Start() {
	// 装载所有GRPC服务
	app.LoadGrpcApp(s.svr)

	// 启动HTTP服务
	lis, err := net.Listen("tcp", s.c.App.GRPC.Addr())
	if err != nil {
		s.l.Errorf("listen grpc tcp conn error, %s", err)
		return
	}

	time.AfterFunc(5*time.Second, s.registry)

	s.l.Infof("GRPC 服务监听地址: %s", s.c.App.GRPC.Addr())
	if err := s.svr.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			s.l.Info("service is stopped")
		}

		s.l.Error("start grpc service error, %s", err.Error())
		return
	}
}

func (s *GRPCService) registry() {
	// 1.获取mcenter sdk实例

	// sdk 提供注册方法

	req := instance.NewRegistryRequest()
	req.Address = s.c.App.GRPC.Addr()
	lf, err := rpc.C().Registry(s.ctx, req)

	if err != nil {
		s.l.Errorf("registry to mcenter error,%s", err)
		return
	}
	// 注销时需要使用
	s.lifecycler = lf
	s.l.Info("registry to mcenter success")
}

// Stop 启动GRPC服务
func (s *GRPCService) Stop() error {
	// 提前剔除注册中心的地址
	if s.lifecycler != nil {
		if err := s.lifecycler.UnRegistry(s.ctx); err != nil {
			s.l.Errorf("unregistry error, %s", err)
		} else {
			s.l.Info("unregistry success")
		}
	}

	s.svr.GracefulStop()
	return nil
}
