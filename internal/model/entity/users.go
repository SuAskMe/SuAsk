// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id           int         `json:"id"           orm:"id"             description:"用户ID"`               // 用户ID
	Name         string      `json:"name"         orm:"name"           description:"用户名"`                // 用户名
	Email        string      `json:"email"        orm:"email"          description:"邮箱"`                 // 邮箱
	PasswordHash string      `json:"passwordHash" orm:"password_hash"  description:"密码哈希，算法暂定为Argon2id"` // 密码哈希，算法暂定为Argon2id
	Role         string      `json:"role"         orm:"role"           description:"角色"`                 // 角色
	Nickname     string      `json:"nickname"     orm:"nickname"       description:"昵称"`                 // 昵称
	Introduction string      `json:"introduction" orm:"introduction"   description:"简介"`                 // 简介
	AvatarFileId int         `json:"avatarFileId" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"` // 头像文件ID，为空时为配置的默认头像
	ThemeId      int         `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`   // 主题ID，为空时为配置的默认主题
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"     description:"创建时间"`               // 创建时间
}
