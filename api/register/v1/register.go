package v1

import "github.com/gogf/gf/v2/frame/g"

type Role string

const (
	Teacher Role = "teacher"
	Student Role = "student"
	Admin   Role = "admin"
)

type RegisterReq struct {
	g.Meta   `path:"/user/register" tags:"Register" method:"POST" summary:"注册接口"`
	Name     string `json:"name" v:"required" dc:"用户名"`
	Password string `json:"password" v:"required" dc:"密码"`
	Role     Role   `json:"role" v:"enums" dc:"角色"`
	Token    string `json:"token" v:"required" dc:"注册邮箱成功时传递的Token，用于在这里验证为同一个人，里面搭载 Email"`
	Email    string `json:"email" v:"required|email" dc:"注册邮箱"`
}

type RegisterRes struct {
	Id int `json:"id" dc:"注册成功的用户ID"`
}
