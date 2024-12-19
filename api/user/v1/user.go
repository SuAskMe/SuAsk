package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UserInfoReq struct {
	g.Meta `path:"/user/info" method:"GET" tags:"User" summary:"请求用户信息"`
}

type UserInfoRes struct {
	UserInfoBase
	Email   string `json:"email"        orm:"email"          description:"邮箱"`
	ThemeId int    `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`
}

type UserInfoByIdReq struct {
	g.Meta `path:"/user/get-user" method:"GET" tags:"User" summary:"通过Id获取用户信息"`
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
	g.Meta        `path:"/user/update-user" method:"POST" tags:"User" summary:"更新用户信息"`
	Nickname      interface{} `json:"nickname"     orm:"nickname"       description:"昵称"`
	Introduction  interface{} `json:"introduction" orm:"introduction"   description:"简介"`
	AvatarFieldId interface{} `json:"avatarId"   description:"头像文件ID，为空时为配置的默认头像"`
	ThemeId       interface{} `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`
}

type UpdateUserRes struct {
	Id int `json:"id"           orm:"id"             description:"用户ID"`
}

type UpdatePasswordReq struct {
	g.Meta   `path:"/user/update-password" method:"POST" tags:"User" summary:"更新密码"`
	Id       string `json:"id" v:"required" dc:"用户ID"`
	Password string `json:"password" v:"required" dc:"新的密码"`
}

type UpdatePasswordRes struct {
	Id int `json:"id"           orm:"id"             description:"用户ID"`
}
