package register

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/golang-jwt/jwt/v5"
	v1 "suask/api/register/v1"
	"suask/internal/model"
	"suask/internal/service"
	"suask/utility"
	"suask/utility/send_code"
	"time"
)

type cRegister struct{}

var Register = cRegister{}

func (c *cRegister) SendVerificationCode(ctx context.Context, req *v1.SendVerificationCodeReq) (res *v1.SendVerificationCodeRes, err error) {
	data := model.CheckEmailAndNameInput{}
	err = gconv.Struct(req, &data)
	if err != nil {
		return nil, err
	}
	out, err := service.Register().CheckEmailAndName(ctx, data)
	if err != nil {
		return nil, err
	}
	var codeTest string
	// 如果都不重复
	if !out.NameDuplicated && !out.EmailDuplicated {
		code, err := send_code.SendCode(data.Email)
		if err != nil {
			return nil, err
		}
		duplicated, _ := gcache.SetIfNotExist(ctx, req.Email, code, 5*time.Minute)
		if !duplicated {
			_, _, err := gcache.Update(ctx, req.Email, code)
			if err != nil {
				return nil, err
			}
			//return nil, gerror.New("你的邮箱和别人重复了")
		}
		codeTest = code
	} else {
		return nil, gerror.New("邮箱或用户名重复")
	}
	return &v1.SendVerificationCodeRes{Code: codeTest}, nil
}

type VerifyClaims struct { // 对邮件验证的token
	Email string
	jwt.RegisteredClaims
}

func (c *cRegister) VerifyVerificationCode(ctx context.Context, req *v1.VerifyVerificationCodeReq) (res *v1.VerifyVerificationCodeRes, err error) {
	code, err := gcache.Get(ctx, req.Email)
	if code == nil {
		return nil, gerror.New("怎么有人偷跑，你压根没获取验证码好吧")
	}
	verificationCode := gconv.String(code)
	if verificationCode != req.Code {
		return nil, gerror.New("验证码错误")
	}
	verifyClaims := &VerifyClaims{
		Email: req.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(6 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, verifyClaims)
	key := utility.JwtKey
	tokenString, err := token.SignedString(key)
	if err != nil {
		return nil, err
	}
	res = &v1.VerifyVerificationCodeRes{}
	res.Token = tokenString
	return res, nil
}

func (c *cRegister) Register(ctx context.Context, req *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	key := utility.JwtKey
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
	out, err := service.Register().Register(ctx, data)
	if err != nil {
		return nil, err
	}
	return &v1.RegisterRes{Id: out.Id}, nil
}
