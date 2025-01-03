package user

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
	"suask/utility/login"
)

type sUser struct {
}

func (s sUser) GetUser(ctx context.Context, in model.UserInfoInput) (out model.UserInfoOutput, err error) {
	user := model.UserInfoOutput{}
	err = dao.Users.Ctx(ctx).Where(do.Users{Id: in.Id}).Scan(&user)
	return user, err
}

func (s sUser) UpdateUser(ctx context.Context, in model.UpdateUserInput) (out model.UpdateUserOutput, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	userInfo := do.Users{
		Nickname:     in.Nickname,
		Introduction: in.Introduction,
		ThemeId:      in.ThemeId,
		AvatarFileId: in.AvatarFileId,
	}
	_, err = dao.Users.Ctx(ctx).WherePri(userId).Update(userInfo)
	if err != nil {
		return model.UpdateUserOutput{}, err
	}
	return model.UpdateUserOutput{Id: userId}, nil
}

func (s sUser) UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error) {
	salt := grand.S(10)
	password := login.EncryptPassword(in.Password, salt)
	update := do.Users{
		Salt:         salt,
		PasswordHash: password,
	}
	_, err = dao.Users.Ctx(ctx).WherePri(in.UserId).Update(update)
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
