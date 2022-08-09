package impl

import (
	"context"
	"fmt"

	"codeup.aliyun.com/baber/go/keyauth/apps/token"
	"github.com/infraboard/mcube/exception"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (i *impl) save(ctx context.Context, ins *token.Token) error {
	if _, err := i.col.InsertOne(ctx, ins); err != nil {
		return exception.NewInternalServerError("inserted token(%s) document error, %s",
			ins.AccessToken, err)
	}
	return nil
}

func (i *impl) get(ctx context.Context, accessToken string) (*token.Token, error) {
	tk, err := i.getFromDB(ctx, accessToken)
	if err != nil {
		return nil, err
	}

	return tk, nil
}

func (i *impl) update(ctx context.Context, ins *token.Token) error {
	if _, err := i.col.UpdateByID(ctx, ins.AccessToken, ins); err != nil {
		return exception.NewInternalServerError("update token(%s) document error, %s", ins.AccessToken, err)
	}

	return nil
}

func (i *impl) deleteToken(ctx context.Context, ins *token.Token) error {
	if ins == nil || ins.AccessToken == "" {
		return fmt.Errorf("token is nil")
	}

	result, err := i.col.DeleteOne(ctx, bson.M{"_id": ins.AccessToken})
	if err != nil {
		return exception.NewInternalServerError("delete token(%s) error, %s", ins.AccessToken, err)
	}

	if result.DeletedCount == 0 {
		return exception.NewNotFound("token %s not found", ins.AccessToken)
	}

	return nil
}

func (s *impl) getFromDB(ctx context.Context, accessToken string) (*token.Token, error) {
	filter := bson.M{"_id": accessToken}

	ins := token.NewDefaultToken()
	if err := s.col.FindOne(ctx, filter).Decode(ins); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, exception.NewNotFound("access token %s not found", accessToken)
		}

		return nil, exception.NewInternalServerError("find access token %s error, %s", accessToken, err)
	}

	return ins, nil
}
