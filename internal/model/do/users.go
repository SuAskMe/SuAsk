// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure of table users for DAO operations like Where/Data.
type Users struct {
	g.Meta       `orm:"table:users, do:true"`
	Id           any         // 用户ID
	Name         any         // 用户名
	Email        any         // 邮箱
	Salt         any         // 加密盐
	PasswordHash any         // 密码哈希，算法暂定为Argon2id
	Role         any         // 角色
	Nickname     any         // 昵称
	Introduction any         // 简介
	AvatarFileId any         // 头像文件ID，为空时为配置的默认头像
	CreatedAt    *gtime.Time // 创建时间
	UpdatedAt    *gtime.Time //
	DeletedAt    *gtime.Time //
}
