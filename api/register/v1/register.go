package v1

import "github.com/gogf/gf/v2/frame/g"

type Role string

const (
	Teacher Role = "teacher"
	Student Role = "student"
	Admin   Role = "admin"
)

type RegisterReq struct {
	g.Meta   `path:"/register" tags:"Register" method:"POST" summary:"注册接口"`
	Name     string `json:"name" v:"required" dc:"用户名"`
	Password string `json:"password" v:"required" dc:"密码"`
	Email    string `json:"email" v:"required|email" dc:"注册邮箱"`
	Code     string `json:"code" v:"required" dc:"邮箱验证码"`
}

type RegisterRes struct {
	Id int `json:"id" dc:"注册成功的用户ID"`
}
