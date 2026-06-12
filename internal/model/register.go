package model

type RegisterInput struct {
	Name     string `json:"name" orm:"name" dc:"用户名"`
	UserSalt string `json:"userSalt" orm:"salt" dc:"加密盐"`
	Password string `json:"password" orm:"password" dc:"密码"`
	Role     string `json:"role" orm:"role" dc:"角色"`
	Email    string `json:"email" orm:"email" dc:"注册邮箱"`
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
