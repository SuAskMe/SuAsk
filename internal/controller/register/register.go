package register

import (
	"context"
	v1 "suask/api/register/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/send_email"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type cRegister struct{}

var Register = cRegister{}

// SendVerificationCode 发送注册验证码（不限邮箱后缀，身份由校园网 IP 中间件保证）。
func (c *cRegister) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error) {
	data := model.CheckEmailAndNameInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}

	// 检查是否已经发送过验证码
	v, err := g.Redis().Get(ctx, consts.RedisSendCodePrefix+req.Email)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	if v.String() != "" {
		return &v1.SendVerificationCodeRes{Msg: "验证码已发送，请注意查收或稍后再试"}, nil
	}

	// 检查邮箱和用户名是否重复
	out, err := service.Register().CheckEmailAndName(ctx, data)
	if err != nil {
		return nil, err
	}
	if out.NameDuplicated || out.EmailDuplicated {
		return &v1.SendVerificationCodeRes{Msg: "邮箱或用户名重复"}, nil
	}

	// 发送验证码
	code, err := send_email.SendCode(data.Email)
	if err != nil {
		return nil, err
	}
	err = g.Redis().SetEX(ctx, consts.RedisSendCodePrefix+req.Email, code, 300)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	err = g.Redis().SetEX(ctx, consts.RedisCountCodePrefix+req.Email, 10, 300)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	return &v1.SendVerificationCodeRes{Msg: "200"}, nil
}

// Register 注册用户：校验验证码 + 创建用户（校园网 IP 校验由中间件完成）。
func (c *cRegister) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	// 校验验证码
	code, err := g.Redis().Get(ctx, consts.RedisSendCodePrefix+req.Email)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	storedCode := code.String()
	if storedCode == "" {
		return nil, gerror.New("验证码已过期，请重新获取")
	}
	if storedCode != req.Code {
		cnt, err := g.Redis().Decr(ctx, consts.RedisCountCodePrefix+req.Email)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, gerror.New(consts.ErrInternal)
		}
		if cnt <= 0 {
			g.Redis().Del(ctx, consts.RedisSendCodePrefix+req.Email, consts.RedisCountCodePrefix+req.Email)
		}
		return nil, gerror.New("验证码错误")
	}
	// 验证码一次性使用
	g.Redis().Del(ctx, consts.RedisSendCodePrefix+req.Email, consts.RedisCountCodePrefix+req.Email)

	// 注册用户
	data := model.RegisterInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	out, err := service.Register().Register(ctx, data)
	if err != nil {
		return nil, err
	}
	// 注册 setting 表，notify_email 默认等于注册邮箱
	_, err = service.Setting().AddSetting(ctx, model.AddSettingInput{
		Id:           out.Id,
		ThemeId:      consts.DefaultThemeId,
		NotifySwitch: true,
		NotifyEmail:  data.Email,
	})
	if err != nil {
		return nil, err
	}
	return &v1.RegisterRes{Id: out.Id}, nil
}
