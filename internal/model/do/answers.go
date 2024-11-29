// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Answers is the golang structure of table answers for DAO operations like Where/Data.
type Answers struct {
	g.Meta     `orm:"table:answers, do:true"`
	Id         interface{} // 回答ID
	UserId     interface{} // 用户ID
	QuestionId interface{} // 问题ID
	InReplyTo  interface{} // 回复的回答ID，可为空
	Contents   interface{} // 回答内容
	CreatedAt  *gtime.Time // 创建时间
	Views      interface{} // 浏览量
}
