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
	userInfo := model.UpdateUserInput{}
	err = gconv.Struct(req, &userInfo)
	if err != nil {
		return nil, err
	}
	// 上传头像
	if req.AvatarFile != nil {
		avatarFile := model.FileUploadInput{File: req.AvatarFile}
		data, err := service.File().UploadFile(ctx, avatarFile)
		if err != nil {
			return nil, err
		}
		avatarId := data.Id
		userInfo.AvatarFileId = avatarId
	}
	out, err := service.User().UpdateUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}
	return &v1.UpdateUserRes{Id: out.Id}, nil
}

func (c *cUser) UpdatePassWord(ctx context.Context, req *v1.UpdatePasswordReq) (res *v1.UpdatePasswordRes, err error) {
	userInfo := model.UpdatePasswordInput{}
	err = gconv.Struct(req, &userInfo)
	if err != nil {
		return nil, err
	}
	out, err := service.User().UpdatePassword(ctx, userInfo)
	if err != nil {
		return nil, err
	}
	res = &v1.UpdatePasswordRes{Id: out.Id}
	return res, nil
}

func (c *cUser) GetUserInfoById(ctx context.Context, req *v1.UserInfoByIdReq) (res *v1.UserInfoByIdRes, err error) {
	UserId := model.UserInfoInput{Id: req.Id}
	out, err := service.User().GetUser(ctx, UserId)
	if err != nil {
		return nil, err
	}
	res = &v1.UserInfoByIdRes{}
	res.Id = out.Id
	res.Name = out.Name
	res.Nickname = out.Nickname
	res.Role = out.Role
	res.Introduction = out.Introduction
	avatarId := out.AvatarFileId
	if avatarId != 0 {
		file, err1 := service.File().Get(ctx, model.FileGetInput{Id: avatarId})
		if err1 != nil {
			return nil, err1
		}
		avatarURL := file.URL
		res.UserInfoBase.AvatarURL = avatarURL
	}
	return res, nil
}

func (c *cUser) Info(ctx context.Context, req *v1.UserInfoReq) (res *v1.UserInfoRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	out, err := service.User().GetUser(ctx, model.UserInfoInput{Id: userId})
	if err != nil {
		return nil, err
	}
	res = &v1.UserInfoRes{}
	res.Id = out.Id
	res.Name = out.Name
	res.Nickname = out.Nickname
	res.Role = out.Role
	res.Introduction = out.Introduction
	avatarId := out.AvatarFileId
	if avatarId != 0 {
		file, err1 := service.File().Get(ctx, model.FileGetInput{Id: avatarId})
		if err1 != nil {
			return nil, err1
		}
		res.AvatarURL = file.URL
	}
	res.Email = out.Email
	res.ThemeId = out.ThemeId
	return res, nil
}
