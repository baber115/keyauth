package token

import (
	"fmt"
	"time"

	"codeup.aliyun.com/baber/go/keyauth/common/utils"
	"github.com/infraboard/mcube/exception"
)

const (
	AppName       = "token"
	DefaultDomain = "domain"
)

func NewIssueTokenRequest() *IssueTokenRequest {
	return &IssueTokenRequest{
		UserDomain: DefaultDomain,
	}
}

func (req *IssueTokenRequest) Validate() error {
	switch req.GranteType {
	case GranteType_PASSWORD:
		if req.UserName == "" || req.Password == "" {
			return fmt.Errorf("password grand required username and password")
		}
	case GranteType_LADP:
	case GranteType_ACCESS_TOKEN:
	case GranteType_REFRSH_TOKEN:
	}
	return nil
}

func NewRevokeTokenRequest() *RevokeTokenRequest {
	return &RevokeTokenRequest{}
}

func NewValidateTokenRequest(AccessToken string) *ValidateTokenRequest {
	return &ValidateTokenRequest{
		AccessToken: AccessToken,
	}
}

func NewQueryTokenRequest() *QueryTokenRequest {
	return &QueryTokenRequest{}
}

func NewToken(req *IssueTokenRequest, duration time.Duration) *Token {
	now := time.Now()
	expired := now.Add(duration)
	refreshExpired := now.Add(duration * 4)

	return &Token{
		AccessToken:           utils.MakeBearer(24),
		IssueAt:               now.UnixMilli(),
		Data:                  req,
		AccessTokenExpiredAt:  expired.UnixMilli(),
		RefreshToken:          utils.MakeBearer(32),
		RefreshTokenExpiredAt: refreshExpired.UnixMilli(),
	}
}

func NewDefaultToken() *Token {
	return &Token{
		Data: &IssueTokenRequest{},
		Meta: map[string]string{},
	}
}

func (t *Token) Validate() error {
	// 判断access token过没过期
	fmt.Println(time.Now().UnixMilli(), t.AccessTokenExpiredAt)
	if time.Now().UnixMilli() > t.AccessTokenExpiredAt {
		return exception.NewAccessTokenExpired("access token expired")
	}

	return nil
}

func (t *Token) IsRefreshTokenExpired() bool {
	// 判断refresh token过没过期
	if time.Now().UnixMilli() > t.RefreshTokenExpiredAt {
		return false
	}

	return true
}

// 续约token，延长一个生命周期
func (t *Token) ExtendToken(expiredDuration time.Duration) {
	now := time.Now()
	t.AccessTokenExpiredAt = now.Add(expiredDuration).UnixMilli()
	t.RefreshTokenExpiredAt = now.Add(expiredDuration * 5).UnixMilli()
}
