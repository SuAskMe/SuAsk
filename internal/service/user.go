package service

import (
	"context"
	"suask/internal/model"
)

type (
	IUser interface {
		UpdateUserInfo(ctx context.Context, in model.UpdateUserInput) (out model.UpdateUserOutput, err error)
		UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
