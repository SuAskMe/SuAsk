// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Teachers is the golang structure of table teachers for DAO operations like Where/Data.
type Teachers struct {
	g.Meta       `orm:"table:teachers, do:true"`
	Id           any //
	Responses    any // 回复数
	Name         any // 老师名字
	AvatarUrl    any // 老师头像链接
	Introduction any // 老师简介
	Email        any // 老师邮箱
	Perm         any // 提问箱权限
}
