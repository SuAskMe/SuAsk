package service

import (
	"context"
	"suask/internal/model"
)

type (
	IRegister interface {
		Register(ctx context.Context, in model.RegisterInput) (out model.RegisterOutput, err error)
		SendVerificationCode(ctx context.Context, in model.SendVerificationCodeInput) (out model.SendVerificationCodeOutput, err error)
		VerifyVerificationCode(ctx context.Context, in model.VerifyVerificationCodeInput) (out model.VerifyVerificationCodeOutput, err error)
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
