package utils_test

import (
	"testing"

	"codeup.aliyun.com/baber/go/keyauth/common/utils"
)

// a9993e364706816aba3e25717850c26c9cd0d89d
func TestHash(t *testing.T) {
	v := utils.Hash("abc")
	t.Log(v)
}

//
// $2a$14$OoDnJEfZZpX0gT9OCYGhIOpyKSigR6rwy4QpgVHA8jLZziAVNqa9C
func TestPasswordHash(t *testing.T) {
	v := utils.HashPassword("abc")
	t.Log(v)
	ok := utils.CheckPasswordHash("abc", "$2a$14$OoDnJEfZZpX0gT9OCYGhIOpyKSigR6rwy4QpgVHA8jLZziAVNqa9C")
	t.Log(ok)
	ok = utils.CheckPasswordHash("abc", "$2a$14$uK2YKOnmjLGV3lHpiXa6LONhvWd2ja9unZP0pEc7XmyeY2li5KwqC")
	t.Log(ok)
}
