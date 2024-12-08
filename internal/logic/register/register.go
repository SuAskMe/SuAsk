package register

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/service"
	"suask/utility"

	"github.com/gogf/gf/v2/util/grand"
)

type sRegister struct{}

func init() {
	service.RegisterRegister(New())
}

func New() *sRegister {
	return &sRegister{}
}

func (s *sRegister) Register(ctx context.Context, in model.RegisterInput) (out model.RegisterOutput, err error) {
	UserSalt := grand.S(10)
	in.Password = utility.EncryptPassword(in.Password, UserSalt)
	in.UserSalt = UserSalt
	in.Role = consts.STUDENT
	lastInsertID, err := dao.Users.Ctx(ctx).InsertAndGetId(in)
	if err != nil {
		return out, err
	}
	return model.RegisterOutput{Id: int(lastInsertID)}, err
}

func (s *sRegister) SendVerificationCode(ctx context.Context, in model.SendVerificationCodeInput) (out model.SendVerificationCodeOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func (s *sRegister) VerifyVerificationCode(ctx context.Context, in model.VerifyVerificationCodeInput) (out model.VerifyVerificationCodeOutput, err error) {
	//TODO implement me
	panic("implement me")
}
