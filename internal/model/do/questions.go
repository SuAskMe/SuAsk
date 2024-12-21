// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure of table questions for DAO operations like Where/Data.
type Questions struct {
	g.Meta    `orm:"table:questions, do:true"`
	Id        interface{} // 问题ID
	SrcUserId interface{} // 发起提问的用户ID
	DstUserId interface{} // 被提问的用户ID，为空时问大家，不为空时问教师
	Title     interface{} // 问题标题
	Contents  interface{} // 问题内容
	IsPrivate interface{} // 是否私密提问，仅在问教师时可为是
	CreatedAt *gtime.Time // 创建时间
	Views     interface{} // 浏览量
}
