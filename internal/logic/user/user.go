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
	"suask/utility"
)

type sUser struct {
}

func (s sUser) GetUserInfoById(ctx context.Context, in model.GetUserInfoByIdInput) (out model.GetUserInfoByIdOutput, err error) {
	user := model.GetUserInfoByIdOutput{}
	err = dao.Users.Ctx(ctx).Where(do.Users{Id: in.Id}).Scan(&user)
	return model.GetUserInfoByIdOutput{UserInfoBase: user.UserInfoBase}, err
}

func (s sUser) GetUserInfo(ctx context.Context, in model.UserInfoInput) (out model.UserInfoOutput, err error) {
	user := model.UserInfoOutput{}
	err = dao.Users.Ctx(ctx).Where(do.Users{Id: in.Id}).Scan(&user)
	return user, err
}

func (s sUser) UpdateUserInfo(ctx context.Context, in model.UpdateUserInput) (out model.UpdateUserOutput, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	_, err = dao.Users.Ctx(ctx).WherePri(userId).Update(in)
	if err != nil {
		return model.UpdateUserOutput{}, err
	}
	return model.UpdateUserOutput{Id: userId}, nil
}

func (s sUser) UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	userSalt := grand.S(10)
	in.Salt = userSalt
	in.Password = utility.EncryptPassword(in.Password, userSalt)
	_, err = dao.Users.Ctx(ctx).WherePri(userId).Update(in)
	if err != nil {
		return model.UpdatePasswordOutput{}, err
	}
	return model.UpdatePasswordOutput{Id: userId}, err
}

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}
