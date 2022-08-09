package rpc

import (
	"context"
	"fmt"
	"time"

	"codeup.aliyun.com/baber/go/keyauth/apps/book"
	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcenter/client/rpc/auth"
	"github.com/infraboard/mcenter/client/rpc/resolver"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	client *ClientSet
)

// SetGlobal todo
func SetGlobal(cli *ClientSet) {
	client = cli
}

// C Global
func C() *ClientSet {
	return client
}

// NewClient todo
// 传递注册中心的地址
func NewClient(conf *rpc.Config) (*ClientSet, error) {
	zap.DevelopmentSetup()
	log := zap.L()
	fmt.Println(conf)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// resolver 进行解析的时候 需要mcenter客户端实例已经初始化
	conn, err := grpc.DialContext(
		ctx,
		fmt.Sprintf("%s://%s", resolver.Scheme, "keyauth"), // Dial to "mcenter://keyauth"
		grpc.WithPerRPCCredentials(auth.NewAuthentication("a", "b")),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &ClientSet{
		conn: conn,
		log:  log,
	}, nil
}

// Client 客户端N
type ClientSet struct {
	conn *grpc.ClientConn
	log  logger.Logger
}

// Token服务的SDK
func (c *ClientSet) Token() token.ServiceClient {
	return token.NewServiceClient(c.conn)
}

// Book服务的SDK
func (c *ClientSet) Book() book.ServiceClient {
	return book.NewServiceClient(c.conn)
}

//// Endpoint服务的SDK
//func (c *ClientSet) Endpoint() endpoint.RPCClient {
//	return endpoint.NewRPCClient(c.conn)
//}
//
//// Role服务的SDK
//func (c *ClientSet) Role() role.RPCClient {
//	return role.NewRPCClient(c.conn)
//}
//
//// Policy服务的SDK
//func (c *ClientSet) Policy() policy.RPCClient {
//	return policy.NewRPCClient(c.conn)
//}
//
//// Audit服务的SDK
//func (c *ClientSet) Audit() audit.RPCClient {
//	return audit.NewRPCClient(c.conn)
//}
