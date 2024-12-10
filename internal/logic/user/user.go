package user

import (
	"context"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
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
	//TODO implement me
	panic("implement me")
}

func (s sUser) UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error) {
	//TODO implement me
	panic("implement me")
}

func init() {
	service.RegisterUser(New())
}

func New() *sUser {
	return &sUser{}
}
