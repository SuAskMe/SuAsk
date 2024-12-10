package user

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "suask/api/user/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"
)

type cUser struct {
}

var User cUser

func (c *cUser) UpdateUserInfo(ctx context.Context, req *v1.UpdateUserReq) (res *v1.UpdateUserRes, err error) {
	return
}

func (c *cUser) UpdatePassWord(ctx context.Context, req *v1.UpdatePasswordReq) (res *v1.UpdatePasswordRes, err error) {
	return
}

func (c *cUser) GetUserInfoById(ctx context.Context, req *v1.UserInfoByIdReq) (res *v1.UserInfoByIdRes, err error) {
	UserId := model.GetUserInfoByIdInput{Id: req.Id}
	out, err := service.User().GetUserInfoById(ctx, UserId)
	if err != nil {
		return nil, err
	}
	return &v1.UserInfoByIdRes{UserInfoBase: out.UserInfoBase}, nil
}

func (c *cUser) Info(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	out, err := service.User().GetUserInfo(ctx, model.UserInfoInput{Id: userId})
	if err != nil {
		return nil, err
	}
	return &v1.UserInfoRes{UserInfoBase: out.UserInfoBase, Email: out.Email, ThemeId: out.ThemeId}, nil
}
