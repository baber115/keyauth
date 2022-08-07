package user

import (
	"net/http"
	"time"

	"codeup.aliyun.com/baber/go/keyauth/common/utils"
	"github.com/go-playground/validator/v10"
	"github.com/infraboard/mcube/http/request"
	pb_request "github.com/infraboard/mcube/pb/request"
	"github.com/rs/xid"
)

const (
	AppName = "user"
)

var (
	validate = validator.New()
)

func (req *CreateUserRequest) Validate() error {
	return validate.Struct(req)
}

// 保存Hash过后的Password
func NewUser(req *CreateUserRequest) *User {
	return &User{
		Id:       xid.New().String(),
		CreateAt: time.Now().UnixMilli(),
		Data:     req,
	}
}
func NewDefaultUser() *User {
	return NewUser(&CreateUserRequest{})
}

func NewUserSet() *UserSet {
	return &UserSet{
		Items: []*User{},
	}
}

func (s *UserSet) Add(item *User) {
	s.Items = append(s.Items, item)
}

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{}
}

func NewDescribeUserRequest(userId string) *DescribeUserRequest {
	return &DescribeUserRequest{
		UserId: userId,
	}
}

func NewDescribeUserRequestByName(domain, name string) *DescribeUserRequest {
	return &DescribeUserRequest{
		DescribeBy: DescribeBy_USER_NAME,
		Domain:     domain,
		UserName:   name,
	}
}

func NewQueryUserRequest() *QueryUserRequest {
	return &QueryUserRequest{
		Page: request.NewDefaultPageRequest(),
	}
}

func NewQueryUserRequestFromHTTP(r *http.Request) *QueryUserRequest {
	qs := r.URL.Query()

	return &QueryUserRequest{
		Page:     request.NewPageRequestFromHTTP(r),
		Keywords: qs.Get("keywords"),
	}
}

func NewPutUserRequest(id string) *UpdateUserRequest {
	return &UpdateUserRequest{
		Id:         id,
		UpdateMode: pb_request.UpdateMode_PUT,
		UpdateAt:   time.Now().UnixMicro(),
		Data:       NewCreateUserRequest(),
	}
}

func NewPatchUserRequest(id string) *UpdateUserRequest {
	return &UpdateUserRequest{
		Id:         id,
		UpdateMode: pb_request.UpdateMode_PATCH,
		UpdateAt:   time.Now().UnixMicro(),
		Data:       NewCreateUserRequest(),
	}
}

func NewDeleteUserRequestWithID(id string) *DeleteUserRequest {
	return &DeleteUserRequest{
		Id: id,
	}
}

func (u *User) CheckPassword(password string) bool {
	return utils.CheckPasswordHash(password, u.Data.Password)
}
