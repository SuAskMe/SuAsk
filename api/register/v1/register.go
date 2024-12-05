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
	UserName string `json:"username" dc:"用户名"`
	Password string `json:"password" dc:"密码"`
	Role     Role   `json:"role" dc:"角色"`
	Token    string `json:"token" dc:"注册邮箱成功时传递的Token，用于在这里验证为同一个人，里面搭载 Email"`
}

type RegisterRes struct {
	Result bool `json:"result" dc:"返回是否成功"`
}
