package v1

import "github.com/gogf/gf/v2/frame/g"

type SendVerificationCodeReq struct {
	g.Meta `path:"/register/send-verification-code" tags:"Register" method:"POST" summary:"发送注册验证码"`
	Email  string `json:"email" v:"required|email#|请输入正确的邮箱格式" dc:"要发送的邮箱地址"`
	Name   string `json:"name" v:"required#请输入用户名" dc:"要注册的用户名"`
}

type SendVerificationCodeRes struct {
	Msg string `json:"msg"`
}
