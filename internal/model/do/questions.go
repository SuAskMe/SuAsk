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
	Id        any         // 问题ID
	SrcUserId any         // 发起提问的用户ID
	DstUserId any         // 被提问的用户ID，为空时问大家，不为空时问教师
	Title     any         // 问题标题
	Contents  any         // 问题内容
	IsPrivate any         // 是否私密提问，仅在问教师时可为是
	CreatedAt *gtime.Time // 创建时间
	Views     any         // 浏览量
	ReplyCnt  any         // 回复数
}
