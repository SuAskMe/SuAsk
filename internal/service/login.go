package service

import (
	"context"
	"suask/internal/model"
)

type (
	ILogin interface {
		Login(ctx context.Context, in model.UserLoginInput) error
		Logout(ctx context.Context) error
	}
)

var (
	localLogin ILogin
)

func Login() ILogin {
	if localLogin == nil {
		panic("implement not found for interface ILogin, forgot register?")
	}
	return localLogin
}

func RegisterLogin(i ILogin) {
	localLogin = i
}
