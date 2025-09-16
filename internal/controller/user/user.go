package user

import (
	"context"
	v1 "suask/api/user/v1"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/utility/send_email"
	"suask/utility/validation"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
)

type cUser struct {
}

var User cUser

func (c *cUser) UpdateUserInfo(ctx context.Context, req *v1.UpdateUserReq) (res *v1.UpdateUserRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	userInfo := model.UpdateUserInput{
		UserId:       userId,
		Nickname:     req.Nickname,
		Introduction: req.Introduction,
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
	// 更新基础数据
	out, err := service.User().UpdateUser(ctx, userInfo)
	if err != nil {
		return nil, err
	}

	// 更新设置
	_, err = service.Setting().UpdateSetting(ctx, model.UpdateSettingInput{
		Id:           userId,
		ThemeId:      req.ThemeId,
		NotifySwitch: req.NotifySwitch,
		NotifyEmail:  req.NotifyEmail,
	})
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
	switch req.Type {
	case consts.ResetPassword:
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
	case consts.ForgetPassword:
		user := entity.Users{}
		var count int
		_ = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Email, req.Email).ScanAndCount(&user, &count, false)
		if count == 0 {
			return nil, gerror.New("邮箱不存在")
		}
		userId = user.Id
	default:
		return nil, err
	}
	isSent, _ := gcache.Contains(ctx, req.Email)
	if isSent {
		return &v1.SendVerificationCodeRes{Msg: "验证码已发送，请注意查收或稍后再试"}, nil
	}
	code, err := send_email.SendCode(req.Email)
	if err != nil {
		return nil, err
	}
	duplicated, _ := gcache.SetIfNotExist(ctx, req.Email, CodeCache{Code: code, UserId: userId}, time.Minute)
	if !duplicated {
		// _, _, err := gcache.Update(ctx, req.Email, CodeCache{Code: code, UserId: userId})
		// if err != nil {
		// 	return nil, err
		// }
		return &v1.SendVerificationCodeRes{Msg: "验证码已发送，请注意查收或稍后再试"}, nil
	}

	return &v1.SendVerificationCodeRes{Msg: "200"}, nil
}

func (c *cUser) UpdatePassWord(ctx context.Context, req *v1.UpdatePasswordReq) (res *v1.UpdatePasswordRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	code, _ := gcache.Get(ctx, req.Email)
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
	code, _ := gcache.Get(ctx, req.Email)
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

	// 获取用户基本信息
	user, err := service.User().GetUser(ctx, model.UserInfoInput{Id: userId})
	if err != nil {
		return nil, err
	}
	res = &v1.UserInfoRes{}
	res.Id = user.Id
	res.Name = user.Name
	res.Nickname = user.Nickname
	res.Role = user.Role
	res.Introduction = user.Introduction
	avatarId := user.AvatarFileId
	if avatarId != 0 {
		file, err1 := service.File().Get(ctx, model.FileGetInput{Id: avatarId})
		if err1 != nil {
			return nil, err1
		}
		res.AvatarURL = file.URL
	}
	res.Email = user.Email

	// 获取设置内容
	setting, err := service.Setting().GetSetting(ctx, model.GetSettingInput{Id: userId})
	if err != nil {
		return nil, err
	}
	res.ThemeId = setting.ThemeId
	res.NotifyEmail = setting.NotifyEmail
	res.NotifySwitch = setting.NotifySwitch

	// 获取提问箱权限（如果是教师）
	perm, _ := validation.IsTeacher(ctx, userId)
	res.QuestionBoxPerm = perm

	return res, nil
}
