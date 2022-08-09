package impl

import (
	"context"
	"fmt"
	"time"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"codeup.aliyun.com/baber/go/keyauth/apps/user"
	"codeup.aliyun.com/baber/go/keyauth/common/utils"
	"github.com/infraboard/mcube/exception"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	AUTH_ERROR           = "user or password not found"
	DefaultTokenDuration = 10 * time.Minute
)

// 颁发token
func (i *impl) IssueToken(ctx context.Context, req *token.IssueTokenRequest) (*token.Token, error) {
	if err := req.Validate(); err != nil {
		return nil, exception.NewBadRequest("validate issue token error,%s", err)
	}

	switch req.GranteType {
	case token.GranteType_PASSWORD:
		// 1.先获取用户对象
		userReq := user.NewDescribeUserRequestByName(req.UserDomain, req.UserName)
		u, err := i.user.DescribeUser(ctx, userReq)
		if err != nil {
			if exception.IsNotFoundError(err) {
				return nil, exception.NewUnauthorized(AUTH_ERROR)
			}
			return nil, err
		}
		// 2.验证用户密码是否正确，
		if checkPassword := u.CheckPassword(req.Password); checkPassword == false {
			return nil, exception.NewUnauthorized(AUTH_ERROR)
		}

		// 3.颁发token
		// 4.rfc:Bearer，请求头，key: Authorization,value：bearer <access_token>
		t := token.NewToken(req, DefaultTokenDuration)

		// 5. 脱敏
		t.Data.Password = ""
		// 6.入库
		if err := i.save(ctx, t); err != nil {
			return nil, err
		}
		return t, nil
	case token.GranteType_LADP:
	case token.GranteType_ACCESS_TOKEN:
	default:
		return nil, fmt.Errorf("grade type %s not implemented", req.GranteType)
	}

	return nil, status.Errorf(codes.Unimplemented, "method IssueToken not implemented")
}

// 撤销token
func (i *impl) RevokeToken(ctx context.Context, req *token.RevokeTokenRequest) (*token.Token, error) {
	// 1. 获取AccessToken
	token, err := i.get(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	// 2. 检查RefreshToken是否匹配
	if token.RefreshToken != req.RefreshToken {
		return nil, exception.NewBadRequest("refresh_token error")
	}
	// 3. 删除
	if err := i.deleteToken(ctx, token); err != nil {
		return nil, err
	}

	return token, nil
}

// 验证token
func (i *impl) ValidateToken(ctx context.Context, req *token.ValidateTokenRequest) (*token.Token, error) {
	// 1. 获取AccessToken
	token, err := i.get(ctx, req.AccessToken)
	if err != nil {
		return nil, err
	}
	// 2. 检查RefreshToken是否匹配
	if err := token.Validate(); err != nil {
		// 如果access token过期
		if utils.IsAccessTokenExpired(err) {
			if token.IsRefreshTokenExpired() {
				return nil, exception.NewRefreshTokenExpired("refresh token expired")
			}
			// 如果refresh token 没过期，刷新过期时间
			token.ExtendToken(DefaultTokenDuration)
			if err := i.update(ctx, token); err != nil {
				return nil, err
			}
			// 返回续约后的token
			return token, nil
		}
		return nil, err
	}

	return token, nil
}

func (i *impl) QueryToken(ctx context.Context, req *token.QueryTokenRequest) (*token.TokenSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryToken not implemented")
}

func (i *impl) DescribeToken(ctx context.Context, req *token.DescribeTokenRequest) (*token.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeToken not implemented")
}
