package v1

import "github.com/gogf/gf/v2/frame/g"

type SendVerificationCodeReq struct {
	g.Meta `path:"/user/send-verification-code" tags:"Register" method:"POST" summary:"发送验证码"`
	Email  string `json:"email" v:"required|email#|请输入正确的邮箱格式" dc:"要发送的邮箱地址"`
	Name   string `json:"name" v:"required#请输入用户名" dc:"要注册的用户名"`
}

type SendVerificationCodeRes struct {
	Code string `json:"code"`
}

type VerifyVerificationCodeReq struct {
	g.Meta `path:"/user/verify-verification-code" tags:"Register" method:"POST" summary:"验证验证码"`
	Email  string `json:"email" v:"required" dc:"要验证的邮箱"`
	Code   string `json:"code" v:"required" dc:"要验证的验证码"`
}

type VerifyVerificationCodeRes struct {
	Token string `json:"token" dc:"验证成功的Token"`
}
