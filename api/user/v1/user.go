package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type UserInfoReq struct {
	g.Meta `path:"/user/info" method:"GET" tags:"User" summary:"请求用户信息"`
}

type UserInfoRes struct {
	UserInfoBase
}

type UserInfoBase struct {
	Id           int    `json:"id"           orm:"id"             description:"用户ID"`               // 用户ID
	Name         string `json:"name"         orm:"name"           description:"用户名"`                // 用户名
	Email        string `json:"email"        orm:"email"          description:"邮箱"`                 // 邮箱
	Nickname     string `json:"nickname"     orm:"nickname"       description:"昵称"`                 // 昵称
	Introduction string `json:"introduction" orm:"introduction"   description:"简介"`                 // 简介
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"` // 头像文件ID，为空时为配置的默认头像
	ThemeId      int    `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`   // 主题ID，为空时为配置的默认主题
}

type UpdateUserReq struct {
	g.Meta       `path:"/user/update-user" method:"POST" tags:"user" summary:"更新用户信息"`
	Id           int    `json:"id"           orm:"id"             description:"用户ID"`
	Nickname     string `json:"nickname"     orm:"nickname"       description:"昵称"`
	Introduction string `json:"introduction" orm:"introduction"   description:"简介"`
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"`
	ThemeId      int    `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`
}

type UpdateUserRes struct{}

type UpdatePasswordReq struct {
	g.Meta   `path:"user/updatePassword" method:"POST" tags:"User" summary:"更新密码"`
	Password string `json:"password" v:"required" dc:"新的密码"`
}

type UpdatePasswordRes struct{}
