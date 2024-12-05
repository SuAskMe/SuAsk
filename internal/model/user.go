package model

type Role string

const (
	Teacher Role = "teacher"
	Student Role = "student"
	Admin   Role = "admin"
)

// RegisterInput 注册输入
type RegisterInput struct {
	UserName string `json:"username" dc:"用户名"`
	Password string `json:"password" dc:"密码"`
	Role     Role   `json:"role" dc:"角色"`
	Token    string `json:"token" dc:"注册邮箱成功时传递的Token，用于在这里验证为同一个人，里面搭载 Email"`
}

// RegisterOutput 注册输出
type RegisterOutput struct {
	UserName string `json:"username"`
}

// SendVerificationCodeInput 发送验证码输入
type SendVerificationCodeInput struct {
	Email string `json:"email" v:"required" dc:"要发送的邮箱地址"`
}

// SendVerificationCodeOutput 发送验证码输出
type SendVerificationCodeOutput struct {
	VerificationCode string `json:"verification_code"`
}

// VerifyVerificationCodeInput 验证验证码输入
type VerifyVerificationCodeInput struct {
	Email            string `json:"email" v:"required" dc:"要验证的邮箱"`
	VerificationCode string `json:"verification_code" v:"required" dc:"要验证的验证码"`
}

// VerifyVerificationCodeOutput 验证验证码输出
type VerifyVerificationCodeOutput struct {
	Token string `json:"token" dc:"验证成功的Token"`
}
