package user

import (
	"context"
	"errors"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
	"suask/utility"

	"github.com/gogf/gf/v2/util/grand"
)

type sUser struct {
}

func (s sUser) GetUser(ctx context.Context, in model.UserInfoInput) (out model.UserInfoOutput, err error) {
	user := model.UserInfoOutput{}
	err = dao.Users.Ctx(ctx).Where(do.Users{Id: in.Id}).Scan(&user)
	return user, err
}

func (s sUser) UpdateUser(ctx context.Context, in model.UpdateUserInput) (out model.UpdateUserOutput, err error) {
	//userId := gconv.Int(ctx.Value(consts.CtxId))
	nn, ok := in.Nickname.(string)
	if ok && (len([]rune(nn))) > 20 {
		return out, errors.New("昵称长度不能超过20")
	}
	userInfo := do.Users{
		Nickname:     in.Nickname,
		Introduction: in.Introduction,
		//ThemeId:      in.ThemeId,
		AvatarFileId: in.AvatarFileId,
	}
	_, err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, in.UserId).
		Update(userInfo)
	if err != nil {
		return model.UpdateUserOutput{}, err
	}
	return model.UpdateUserOutput{Id: in.UserId}, nil
}

func (s sUser) UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error) {
	salt := grand.S(10)
	password := utility.EncryptPassword(in.Password, salt)
	update := do.Users{
		Salt:         salt,
		PasswordHash: password,
	}
	md := dao.Users.Ctx(ctx)
	switch in.Type {
	case consts.ResetPassword:
		md = md.Where(dao.Users.Columns().Id, in.UserId)
	case consts.ForgetPassword:
		md = md.Where(dao.Users.Columns().Email, in.Email)
	}
	_, err = md.Update(update)
	if err != nil {
		return model.UpdatePasswordOutput{}, err
	}
	return model.UpdatePasswordOutput{Id: in.UserId}, err
}

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}
