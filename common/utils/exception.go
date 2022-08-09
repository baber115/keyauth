package utils

import "github.com/infraboard/mcube/exception"

// 判断token是否过期
func IsAccessTokenExpired(err error) bool {
	if err == nil {
		return false
	}

	e, ok := err.(exception.APIException)
	if !ok {
		return false
	}

	return e.ErrorCode() == exception.AccessTokenExpired && e.Namespace() == exception.GlobalNamespace.String()
}
