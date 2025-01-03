package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type UserInfoReq struct {
	g.Meta `path:"/user" method:"GET" tags:"User" summary:"请求用户信息"`
}

type UserInfoRes struct {
	UserInfoBase
	Email   string `json:"email"        orm:"email"          description:"邮箱"`
	ThemeId int    `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`
}

type UserInfoByIdReq struct {
	g.Meta `path:"/info/user" method:"GET" tags:"Info" summary:"通过Id获取用户信息"`
	Id     int `json:"id" v:"required" dc:"用户ID"`
}

type UserInfoByIdRes struct {
	UserInfoBase
}

type UserInfoBase struct {
	Id           int    `json:"id"           orm:"id"             description:"用户ID"`
	Name         string `json:"name"         orm:"name"           description:"用户名"`
	Role         string `json:"role"         orm:"role"           description:"角色"`
	Nickname     string `json:"nickname"     orm:"nickname"       description:"昵称"`
	Introduction string `json:"introduction" orm:"introduction"   description:"简介"`
	AvatarURL    string `json:"avatar"       description:"头像文件链接"`
}

type UpdateUserReq struct {
	g.Meta       `path:"/user" method:"PUT" tags:"User" summary:"更新用户信息"`
	Nickname     interface{}       `json:"nickname"     orm:"nickname"       description:"昵称"`
	Introduction interface{}       `json:"introduction" orm:"introduction"   description:"简介"`
	AvatarFile   *ghttp.UploadFile `json:"avatar"       description:"头像文件"`
	ThemeId      interface{}       `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`
}

type UpdateUserRes struct {
	Id int `json:"id"           orm:"id"             description:"用户ID"`
}

type UpdatePasswordReq struct {
	g.Meta   `path:"/user/password" method:"PUT" tags:"User" summary:"更新密码"`
	Email    string `json:"email" v:"required|email" dc:"邮箱"`
	Code     string `json:"code" v:"required" dc:"验证码"`
	Password string `json:"password" v:"required" dc:"新的密码"`
}

type UpdatePasswordRes struct {
	Id int `json:"id"           orm:"id"             description:"用户ID"`
}

type SendVerificationCodeReq struct {
	g.Meta `path:"/user/send-code" method:"POST" tags:"User" summary:"发送验证码"`
	Email  string `json:"email" v:"required|email" dc:"邮箱"`
	Type   string `json:"type" v:"required" dc:"方式"`
}

type SendVerificationCodeRes struct {
	Code string `json:"code" dc:"验证码"`
}

type ForgetPasswordReq struct {
	g.Meta   `path:"/user/forget-password" method:"POST" tags:"User" summary:"忘记密码"`
	Email    string `json:"email" v:"required|email" dc:"邮箱"`
	Code     string `json:"code" v:"required" dc:"验证码"`
	Password string `json:"password" v:"required" dc:"新的密码"`
}

type ForgetPasswordRes struct {
	Id int `json:"id"           orm:"id"             description:"用户ID"`
}
