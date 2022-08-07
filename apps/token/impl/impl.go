package impl

import (
	"context"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"codeup.aliyun.com/baber/go/keyauth/apps/user"
	"codeup.aliyun.com/baber/go/keyauth/conf"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"google.golang.org/grpc"
)

var (
	svr = &impl{}
)

type impl struct {
	col *mongo.Collection
	log logger.Logger
	token.UnimplementedServiceServer

	user user.ServiceServer
}

func (i *impl) Config() error {
	// 依赖MongoDB的DB对象
	db, err := conf.C().Mongo.GetDB()
	if err != nil {
		return err
	}
	// 获取一个Collection对象, 通过Collection对象 来进行CRUD
	i.col = db.Collection(i.Name())
	i.log = zap.L().Named(i.Name())
	i.user = app.GetGrpcApp(user.AppName).(user.ServiceServer)

	// 创建索引
	indexs := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{
				{Key: "refresh_token", Value: bsonx.Int32(-1)},
			},
			Options: options.Index().SetUnique(true),
		},
	}

	_, err = i.col.Indexes().CreateMany(context.Background(), indexs)
	if err != nil {
		return err
	}

	return nil
}

func (i *impl) Name() string {
	return token.AppName
}

func (i *impl) Registry(server *grpc.Server) {
	token.RegisterServiceServer(server, svr)
	user.RegisterServiceServer(server, svr.user)
}

func init() {
	app.RegistryGrpcApp(svr)
}
