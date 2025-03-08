package user

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "suask/api/user/v1"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/utility/send_code"
	"suask/utility/validation"
	"time"
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
	if req.AvatarFile.FileHeader != nil {
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

type CodeCache struct {
	Code   string
	UserId int
}

func (c *cUser) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error) {
	var userId int
	if req.Type == consts.ResetPassword {
		userId = gconv.Int(ctx.Value(consts.CtxId))
		if userId == consts.DefaultUserId {
			return nil, gerror.New("默认用户不能改密码")
		}
		user, err := service.User().GetUser(ctx, model.UserInfoInput{Id: userId})
		if err != nil {
			return nil, err
		}
		if req.Email != user.Email {
			return nil, gerror.New("用户注册邮箱与输入邮箱不同")
		}
	} else if req.Type == consts.ForgetPassword {
		user := entity.Users{}
		var count int
		_ = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Email, req.Email).ScanAndCount(&user, &count, false)
		if count == 0 {
			return nil, gerror.New("邮箱不存在")
		}
		userId = user.Id
	} else {
		return nil, err
	}
	code, err := send_code.SendCode(req.Email)
	if err != nil {
		return nil, err
	}
	duplicated, _ := gcache.SetIfNotExist(ctx, req.Email, CodeCache{Code: code, UserId: userId}, 5*time.Minute)
	if !duplicated {
		_, _, err := gcache.Update(ctx, req.Email, CodeCache{Code: code, UserId: userId})
		if err != nil {
			return nil, err
		}
	}
	res = &v1.SendVerificationCodeRes{
		Code: code,
	}
	return res, nil
}

func (c *cUser) UpdatePassWord(ctx context.Context, req *v1.UpdatePasswordReq) (res *v1.UpdatePasswordRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	code, err := gcache.Get(ctx, req.Email)
	var codeStruct CodeCache
	err = gconv.Scan(code, &codeStruct)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, gerror.New("怎么有人偷跑，你压根没获取验证码好吧")
	}
	verificationCode := gconv.String(codeStruct.Code)
	if verificationCode != req.Code {
		return nil, gerror.New("验证码错误")
	}
	input := model.UpdatePasswordInput{Password: req.Password, UserId: userId}
	out, err := service.User().UpdatePassword(ctx, input)
	if err != nil {
		return nil, err
	}
	res = &v1.UpdatePasswordRes{Id: out.Id}
	return res, nil
}

func (c *cUser) ForgetPassword(ctx context.Context, req *v1.ForgetPasswordReq) (res *v1.ForgetPasswordRes, err error) {
	code, err := gcache.Get(ctx, req.Email)
	var codeStruct CodeCache
	err = gconv.Scan(code, &codeStruct)
	if err != nil {
		return nil, err
	}
	if code == nil {
		return nil, gerror.New("怎么有人偷跑，你压根没获取验证码好吧")
	}
	verificationCode := gconv.String(codeStruct.Code)
	if verificationCode != req.Code {
		return nil, gerror.New("验证码错误")
	}
	input := model.UpdatePasswordInput{Password: req.Password, UserId: codeStruct.UserId}
	out, err := service.User().UpdatePassword(ctx, input)
	if err != nil {
		return nil, err
	}
	res = &v1.ForgetPasswordRes{Id: out.Id}
	return res, nil
}

func (c *cUser) GetUserInfoById(ctx context.Context, req *v1.UserInfoByIdReq) (res *v1.UserInfoByIdRes, err error) {
	UserId := model.UserInfoInput{Id: req.Id}
	out, err := service.User().GetUser(ctx, UserId)
	if err != nil {
		return nil, gerror.New("没有该用户")
	}
	res = &v1.UserInfoByIdRes{}
	res.Id = out.Id
	res.Name = out.Name
	res.Nickname = out.Nickname
	res.Role = out.Role
	res.Introduction = out.Introduction
	avatarId := out.AvatarFileId
	if avatarId != 0 {
		file, err := service.File().Get(ctx, model.FileGetInput{Id: avatarId})
		if err != nil {
			return nil, err
		}
		avatarURL := file.URL
		res.UserInfoBase.AvatarURL = avatarURL
	} else {
		res.AvatarURL = consts.DefaultAvatarURL
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

	perm, _ := validation.IsTeacher(ctx, userId)
	res.QuestionBoxPerm = perm

	return res, nil
}
