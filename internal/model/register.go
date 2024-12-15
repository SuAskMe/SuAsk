package model

type RegisterInput struct {
	Name     string `json:"name" orm:"name" dc:"用户名"`
	UserSalt string `json:"userSalt" orm:"salt" dc:"加密盐"`
	Password string `json:"password" orm:"password" dc:"密码"`
	Role     string `json:"role" orm:"role" dc:"角色"`
	Email    string `json:"email" orm:"email" dc:"注册邮箱"`
	Token    string `json:"token" dc:"注册邮箱成功时传递的Token，用于在这里验证为同一个人，里面搭载 Email"`
}

type RegisterOutput struct {
	Id int `json:"id"`
}

type CheckEmailAndNameInput struct {
	Email string `json:"email" v:"required" dc:"要发送的邮箱地址"`
	Name  string `json:"name" v:"required" dc:"注册的用户名"`
}

type CheckEmailAndNameOutput struct {
	EmailDuplicated bool `dc:"要检查的邮箱"`
	NameDuplicated  bool `dc:"要检查的用户名"`
}

type VerifyCodeInput struct {
	Email            string `json:"email" v:"required" dc:"要验证的邮箱"`
	VerificationCode string `json:"verification_code" v:"required" dc:"要验证的验证码"`
}

type VerifyCodeOutput struct {
	Token string `json:"token" dc:"验证成功的Token"`
}
