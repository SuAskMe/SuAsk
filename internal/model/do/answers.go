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
	Id         any         // 回答ID
	UserId     any         // 用户ID
	QuestionId any         // 问题ID
	InReplyTo  any         // 回复的回答ID，可为空
	Contents   any         // 回答内容
	CreatedAt  *gtime.Time // 创建时间
	Upvotes    any         // 点赞量
}
