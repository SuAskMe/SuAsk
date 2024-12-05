package register

import (
	"context"
	"suask/internal/model"
	"suask/internal/service"
)

type sRegister struct{}

func (s sRegister) Register(ctx context.Context, in model.RegisterInput) (out model.RegisterOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sRegister) SendVerificationCode(ctx context.Context, in model.SendVerificationCodeInput) (out model.SendVerificationCodeOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func (s sRegister) VerifyVerificationCode(ctx context.Context, in model.VerifyVerificationCodeInput) (out model.VerifyVerificationCodeOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	service.RegisterRegister(New())
}

func New() *sRegister {
	return &sRegister{}
}
