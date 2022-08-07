package utils_test

import (
	"testing"

	"codeup.aliyun.com/baber/go/keyauth/common/utils"
)

func TestToken(t *testing.T) {
	v := utils.MakeBearer(24)
	t.Log(v)
}
