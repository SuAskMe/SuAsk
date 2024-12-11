package register

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
	"suask/utility/login"

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
	in.Password = login.EncryptPassword(in.Password, UserSalt)
	in.UserSalt = UserSalt
	in.Role = consts.STUDENT

	registerUser := do.Users{
		Name:         in.Name,
		Email:        in.Token,
		Salt:         UserSalt,
		PasswordHash: in.Password,
		Role:         in.Role,
		Nickname:     in.Name,
		Introduction: "",
		ThemeId:      consts.DefaultThemeId,
	}

	lastInsertID, err := dao.Users.Ctx(ctx).InsertAndGetId(registerUser)
	if err != nil {
		if gstr.Contains(err.Error(), "users.name") {
			return out, gerror.New("用户名重复")
		} else if gstr.Contains(err.Error(), "users.email") {
			return out, gerror.New("邮箱重复")
		}
		return out, err
	}
	return model.RegisterOutput{Id: int(lastInsertID)}, err
}

func (s *sRegister) CheckEmailAndName(ctx context.Context, in model.CheckEmailAndNameInput) (out model.CheckEmailAndNameOutput, err error) {
	out = model.CheckEmailAndNameOutput{}
	nameCount, err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Name, in.Name).Count()
	if err != nil {
		return out, err
	}
	emailCount, err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Email, in.Email).Count()
	if err != nil {
		return out, err
	}
	out.NameDuplicated = false
	out.EmailDuplicated = false
	if nameCount > 0 {
		out.NameDuplicated = true
	}
	if emailCount > 0 {
		out.EmailDuplicated = true
	}
	return out, nil
}
