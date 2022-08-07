package impl_test

import (
	"context"
	"testing"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"codeup.aliyun.com/baber/go/keyauth/conf"
	"github.com/infraboard/mcube/app"
	"github.com/infraboard/mcube/logger/zap"
)

var (
	ins token.ServiceServer
)

func init() {
	if err := conf.LoadConfigFromEnv(); err != nil {
		panic(err)
	}

	zap.DevelopmentSetup()
	if err := app.InitAllApp(); err != nil {
		panic(err)
	}

	ins = app.GetGrpcApp(token.AppName).(token.ServiceServer)
}

func TestIssueToken(t *testing.T) {
	req := token.NewIssueTokenRequest()
	req.UserDomain = "a"
	req.UserName = "bobo"
	req.Password = "123456"
	token, err := ins.IssueToken(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(token)
}
