package impl_test

import (
	"context"
	"testing"

	"codeup.aliyun.com/baber/go/keyauth/apps/user"
	"codeup.aliyun.com/baber/go/keyauth/conf"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	ins user.ServiceServer
)

func init() {
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	zap.DevelopmentSetup()
	if err := app.InitAllApp(); err != nil {
		panic(err)
	}

	ins = app.GetGrpcApp(user.AppName).(user.ServiceServer)
}

func TestUserCreate(t *testing.T) {
	req := user.NewCreateUserRequest()
	req.Name = "123"
	req.Password = "123456"
	req.Domain = "1"
	user, err := ins.CreateUser(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}

func TestUserQuery(t *testing.T) {
	req := user.NewQueryUserRequest()
	user, err := ins.QueryUser(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(user)
}
