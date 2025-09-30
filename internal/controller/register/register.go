package register

import (
	"context"
	"strings"
	v1 "suask/api/register/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/send_email"
	"suask/module/sjwt"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
)

type cRegister struct{}

var Register = cRegister{}

func (c *cRegister) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error) {
	data := model.CheckEmailAndNameInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}

	// 检查是否为学校邮箱（暂时）
	if !strings.HasSuffix(req.Email, "@mail.sysu.edu.cn") && !strings.HasSuffix(req.Email, "@mail2.sysu.edu.cn") {
		return nil, gerror.New("暂不支持非中大邮箱注册")
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
	// 如果都不重复
	if !out.NameDuplicated && !out.EmailDuplicated {
		code, err := send_email.SendCode(data.Email)
		if err != nil {
			return nil, err
		}
		err = g.Redis().SetEX(ctx, consts.RedisSendCodePrefix+req.Email, code, 60)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, gerror.New(consts.ErrInternal)
		}
		err = g.Redis().SetEX(ctx, consts.RedisCountCodePrefix+req.Email, 10, 60)
		if err != nil {
			g.Log().Error(ctx, err)
			return nil, gerror.New(consts.ErrInternal)
		}
	} else {
		return &v1.SendVerificationCodeRes{Msg: "邮箱或用户名重复"}, nil
	}
	return &v1.SendVerificationCodeRes{Msg: "200"}, nil
}

type VerifyClaims struct { // 对邮件验证的token
	Email string
	jwt.RegisteredClaims
}

func (c *cRegister) VerifyVerificationCode(ctx context.Context, req *v1.VerifyVerificationCodeReq) (res *v1.VerifyVerificationCodeRes, err error) {
	code, err := g.Redis().Get(ctx, consts.RedisSendCodePrefix+req.Email)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	verificationCode := code.String()
	if verificationCode == "" {
		return nil, gerror.New("验证码已过期，请重新获取")
	}
	if verificationCode != req.Code {
		// 防止爆破
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
	verifyClaims := &VerifyClaims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, verifyClaims)
	key := sjwt.GetKey()
	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	res = &v1.VerifyVerificationCodeRes{}
	res.Token = tokenString
	return res, nil
}

func (c *cRegister) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	key := sjwt.GetKey()
	tokenClaims, _ := jwt.ParseWithClaims(req.Token, &VerifyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if tokenClaims == nil {
		return nil, gerror.New("Token 错误")
	}
	var tokenEmail string
	if claims, ok := tokenClaims.Claims.(*VerifyClaims); ok && tokenClaims.Valid {
		tokenEmail = claims.Email
	}
	if tokenEmail != req.Email {
		return nil, gerror.New("输入的邮箱和验证的邮箱不一致")
	}
	data := model.RegisterInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	// 注册user表
	out, err := service.Register().Register(ctx, data)
	if err != nil {
		return nil, err
	}
	// 注册setting表
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
