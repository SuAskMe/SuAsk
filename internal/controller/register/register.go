package register

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "suask/api/register/v1"
	"suask/internal/model"
	"suask/internal/service"
)

type cRegister struct{}

var Register = cRegister{}

func (c *cRegister) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	data := model.RegisterInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	out, err := service.Register().Register(ctx, data)
	if err != nil {
		return nil, err
	}
	return &v1.RegisterRes{Id: out.Id}, nil
}

func (c *cRegister) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error) {
	return
}

func (c *cRegister) VerifyVerificationCode(ctx context.Context, req *v1.VerifyVerificationCodeReq) (res *v1.VerifyVerificationCodeRes, err error) {
	return
}
