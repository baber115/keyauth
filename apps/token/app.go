package token

import (
	"fmt"
	"time"

	"codeup.aliyun.com/baber/go/keyauth/common/utils"
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

func NewRevolkTokenRequest() *RevolkTokenRequest {
	return &RevolkTokenRequest{}
}

func NewValidateTokenRequest() *ValidateTokenRequest {
	return &ValidateTokenRequest{}
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
