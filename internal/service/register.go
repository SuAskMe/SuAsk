package service

import (
	"context"
	"suask/internal/model"
)

type (
	IRegister interface {
		Register(ctx context.Context, in model.RegisterInput) (out model.RegisterOutput, err error)
		CheckEmailAndName(ctx context.Context, in model.CheckEmailAndNameInput) (out model.CheckEmailAndNameOutput, err error)
	}
)

var (
	localRegister IRegister
)

func Register() IRegister {
	if localRegister == nil {
		panic("implement not found for interface IRegister, forgot register?")
	}
	return localRegister
}

func RegisterRegister(i IRegister) {
	localRegister = i
}
