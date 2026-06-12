package v1

import "github.com/gogf/gf/v2/frame/g"

// --- 创建临时用户 ---

type GuestLoginReq struct {
	g.Meta `path:"/guest/login" method:"POST" tags:"Guest" summary:"创建临时用户"`
}

type GuestLoginRes struct {
	Role string `json:"role" dc:"用户角色"`
	Id   int    `json:"id"   dc:"用户ID"`
}

// --- 升级为正式用户 ---

type GuestUpgradeReq struct {
	g.Meta   `path:"/guest/upgrade" method:"POST" tags:"Guest" summary:"临时用户升级"`
	Name     string `json:"name"     v:"required|length:1,50" dc:"用户名"`
	Email    string `json:"email"    v:"required|email" dc:"邮箱"`
	Password string `json:"password" v:"required|min-length:6" dc:"密码"`
	Code     string `json:"code"     v:"required" dc:"邮箱验证码"`
}

type GuestUpgradeRes struct {
	Role string `json:"role" dc:"新角色"`
	Id   int    `json:"id"   dc:"用户ID"`
}

// --- 升级时发送验证码 ---

type GuestSendCodeReq struct {
	g.Meta `path:"/guest/send-code" method:"POST" tags:"Guest" summary:"升级时发送验证码"`
	Name   string `json:"name"  v:"required" dc:"要注册的用户名"`
	Email  string `json:"email" v:"required|email" dc:"邮箱"`
}

type GuestSendCodeRes struct {
	Msg string `json:"msg" dc:"消息"`
}
