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
	Id           interface{} //
	Responses    interface{} // 回复数
	Name         interface{} // 老师名字
	AvatarUrl    interface{} // 老师头像链接
	Introduction interface{} // 老师简介
	Email        interface{} // 老师邮箱
	Perm         interface{} // 提问箱权限
}
