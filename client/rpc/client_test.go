package rpc_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"codeup.aliyun.com/baber/go/keyauth/client/rpc"
	mcenter "github.com/infraboard/mcenter/client/rpc"
	"github.com/infraboard/mcenter/client/rpc/auth"
	"github.com/infraboard/mcenter/client/rpc/resolver"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// keyauth 客户端
// 需要配置注册中心的地址
// 获取注册中心的客户端，使用注册中心的客户端 查询 keyauth的地址
func TestBookQuery(t *testing.T) {
	should := assert.New(t)

	conf := mcenter.NewDefaultConfig()
	conf.Address = "127.0.0.1:18010"
	conf.ClientID = "uIcoyPmEN0MxkodNLNb9lgBZ"
	conf.ClientSecret = "f5nk9itz0vzhB6L4jyfJ5pV8h1AlpEC0"
	// 传递Mcenter配置, 客户端通过Mcenter进行搜索
	c, err := rpc.NewClient(conf)

	if should.NoError(err) {
		resp, err := c.Token().ValidateToken(
			context.Background(),
			token.NewValidateTokenRequest("Tz1U6QpnZjdd7FKFXxosntOL"),
		)
		should.NoError(err)
		fmt.Println(resp)
	}
}

func TestConnection(t *testing.T) {
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
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
}

func init() {
	// 提前加载好 mcenter客户端, resolver需要使用
	err := mcenter.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
}
