package user

import (
	"context"
	"suask/internal/model"
	"suask/internal/service"
)

type sUser struct {
}

func (s sUser) UpdateUserInfo(ctx context.Context, in model.UpdateUserInput) (out model.UpdateUserOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sUser) UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}
