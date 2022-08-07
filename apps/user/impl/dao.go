package impl

import (
	"context"
	"fmt"

	"codeup.aliyun.com/baber/go/keyauth/apps/user"
	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newQueryRequest(r *user.QueryUserRequest) *queryUserRequest {
	return &queryUserRequest{
		r,
	}
}

type queryUserRequest struct {
	*user.QueryUserRequest
}

func (r *queryUserRequest) FindOptions() *options.FindOptions {
	pageSize := int64(r.Page.PageSize)
	skip := int64(r.Page.PageSize) * int64(r.Page.PageNumber-1)

	opt := &options.FindOptions{
		Sort: bson.D{
			{Key: "create_at", Value: -1},
		},
		Limit: &pageSize,
		Skip:  &skip,
	}

	return opt
}

func (r *queryUserRequest) FindFilter() bson.M {
	filter := bson.M{}
	if r.Keywords != "" {
		filter["$or"] = bson.A{
			bson.M{"data.name": bson.M{"$regex": r.Keywords, "$options": "im"}},
			bson.M{"data.author": bson.M{"$regex": r.Keywords, "$options": "im"}},
		}
	}
	return filter
}

func (i *impl) get(ctx context.Context, req *user.DescribeUserRequest) (*user.User, error) {
	filter := bson.M{}
	switch req.DescribeBy {
	case user.DescribeBy_USER_ID:
		filter["_id"] = req.UserId
	case user.DescribeBy_USER_NAME:
		filter["data.domain"] = req.Domain
		filter["data.name"] = req.UserName
	default:
		return nil, fmt.Errorf("unknow describe_by %s", req.DescribeBy)
	}

	ins := user.NewDefaultUser()
	if err := i.col.FindOne(ctx, filter).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("user %s not found", req)
		}
		return nil, exception.NewInternalServerError("find user %s error, %s", req, err)
	}

	return ins, nil
}

func (i *impl) query(ctx context.Context, req *queryUserRequest) (*user.UserSet, error) {
	resp, err := i.col.Find(ctx, req.FindFilter(), req.FindOptions())

	if err != nil {
		return nil, exception.NewInternalServerError("find user error, error is %s", err)
	}

	set := user.NewUserSet()
	// 循环
	for resp.Next(ctx) {
		ins := user.NewDefaultUser()
		if err := resp.Decode(ins); err != nil {
			return nil, exception.NewInternalServerError("decode user error, error is %s", err)
		}

		set.Add(ins)
	}

	// count
	count, err := i.col.CountDocuments(ctx, req.FindFilter())
	if err != nil {
		return nil, exception.NewInternalServerError("get user count error, error is %s", err)
	}
	set.Total = count

	return set, nil
}

func (i *impl) update(ctx context.Context, ins *user.User) error {
	if _, err := i.col.UpdateByID(ctx, ins.Id, ins); err != nil {
		return exception.NewInternalServerError("inserted user(%s) document error, %s",
			ins.Data.Name, err)
	}

	return nil
}

func (i *impl) deleteUser(ctx context.Context, ins *user.User) error {
	if ins == nil || ins.Id == "" {
		return fmt.Errorf("user is nil")
	}

	result, err := i.col.DeleteOne(ctx, bson.M{"_id": ins.Id})
	if err != nil {
		return exception.NewInternalServerError("delete user(%s) error, %s", ins.Id, err)
	}

	if result.DeletedCount == 0 {
		return exception.NewNotFound("user %s not found", ins.Id)
	}

	return nil
}
