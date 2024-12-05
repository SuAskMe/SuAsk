package v1

import "github.com/gogf/gf/v2/frame/g"

type SendVerificationCodeReq struct {
	g.Meta `path:"/user/send-verification-code" tags:"VerificationCode" method:"POST" summary:"发送验证码"`
	Email  string `json:"email" v:"required" dc:"要发送的邮箱地址"`
}

type SendVerificationCodeRes struct {
	VerificationCode string `json:"verification_code"`
}

type VerifyVerificationCodeReq struct {
	g.Meta           `path:"/user/verify-verification-code" tags:"VerificationCode" method:"POST" summary:"验证验证码"`
	Email            string `json:"email" v:"required" dc:"要验证的邮箱"`
	VerificationCode string `json:"verification_code" v:"required" dc:"要验证的验证码"`
}

type VerifyVerificationCodeRes struct {
	Token string `json:"token" dc:"验证成功的Token"`
}
