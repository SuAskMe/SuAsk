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
	Id         interface{} // 点赞ID
	UserId     interface{} // 用户ID
	QuestionId interface{} // 问题ID
	AnswerId   interface{} // 回复ID
}
