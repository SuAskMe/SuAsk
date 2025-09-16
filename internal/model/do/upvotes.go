// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Upvotes is the golang structure of table upvotes for DAO operations like Where/Data.
type Upvotes struct {
	g.Meta     `orm:"table:upvotes, do:true"`
	Id         any // 点赞ID
	UserId     any // 用户ID
	QuestionId any // 问题ID
	AnswerId   any // 回复ID
}
