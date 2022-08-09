package rpc_test

import (
	"context"
	"fmt"
	"testing"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"codeup.aliyun.com/baber/go/keyauth/client/rpc"
	"github.com/stretchr/testify/assert"

	mcenter "github.com/infraboard/mcenter/client/rpc"
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

func init() {
	// 提前加载好 mcenter客户端, resolver需要使用
	err := mcenter.LoadClientFromEnv()
	if err != nil {
		panic(err)
	}
}
