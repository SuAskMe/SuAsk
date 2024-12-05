package register

import (
	"context"
	v1 "suask/api/register/v1"
)

type cRegister struct{}

var Register cRegister

func (c *cRegister) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	return
}

func (c *cRegister) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error) {
	return
}

func (c *cRegister) VerifyVerificationCode(ctx context.Context, req *v1.VerifyVerificationCodeReq) (res *v1.VerifyVerificationCodeRes, err error) {
	return
}
