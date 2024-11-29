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
	Id           interface{} // 用户ID
	Name         interface{} // 用户名
	Email        interface{} // 邮箱
	PasswordHash interface{} // 密码哈希，算法暂定为Argon2id
	Role         interface{} // 角色
	Nickname     interface{} // 昵称
	Introduction interface{} // 简介
	AvatarFileId interface{} // 头像文件ID，为空时为配置的默认头像
	ThemeId      interface{} // 主题ID，为空时为配置的默认主题
	CreatedAt    *gtime.Time // 创建时间
}
