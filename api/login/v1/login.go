package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// 登录

type LoginReq struct {
	g.Meta   `path:"/login" method:"POST" tag:"Login" summary:"登录请求"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRes struct {
	Type         string `json:"type"         dc:"Token格式"`
	Token        string `json:"token"        dc:"用户的Token"`
	Role         string `json:"role"         orm:"role"           dc:"用户角色"`
	Id           int    `json:"id"           orm:"id"             dc:"用户ID"`
	Name         string `json:"name"         orm:"name"           dc:"用户名"`
	Email        string `json:"email"        orm:"email"          dc:"邮箱"`
	Nickname     string `json:"nickname"     orm:"nickname"       dc:"昵称"`
	AvatarURL    string `json:"avatar"       dc:"头像URL"`
	ThemeId      int    `json:"themeId"      orm:"theme_id"       dc:"主题ID，为空时为配置的默认主题"`
	Introduction string `json:"introduction" orm:"introduction"   dc:"简介"`
}

// 登出

type LogoutReq struct {
	g.Meta `path:"/user/logout" method:"POST" tag:"Logout" summary:"登出请求"`
}

type LogoutRes struct{}

// 刷新 Token

type RefreshTokenReq struct {
	Token string `json:"token"`
}

type RefreshTokenRes struct {
	Token string `json:"token"`
}
