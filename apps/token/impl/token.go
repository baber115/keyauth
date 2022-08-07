package impl

import (
	"context"
	"fmt"
	"time"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"codeup.aliyun.com/baber/go/keyauth/apps/user"
	"github.com/infraboard/mcube/exception"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	AUTH_ERROR = "user or password not found"
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
		t := token.NewToken(req, 10*time.Minute)
		return t, nil
	case token.GranteType_LADP:
	case token.GranteType_ACCESS_TOKEN:
	case token.GranteType_REFRSH_TOKEN:
	default:
		return nil, fmt.Errorf("grade type %s not implemented", req.GranteType)
	}

	return nil, status.Errorf(codes.Unimplemented, "method IssueToken not implemented")
}

// 撤销token
func (i *impl) RevolkToken(ctx context.Context, req *token.RevolkTokenRequest) (*token.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevolkToken not implemented")
}

// 验证token
func (i *impl) ValidateToken(ctx context.Context, req *token.ValidateTokenRequest) (*token.Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateToken not implemented")
}

func (i *impl) QueryToken(ctx context.Context, req *token.QueryTokenRequest) (*token.TokenSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryToken not implemented")
}
