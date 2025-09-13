// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Notifications is the golang structure of table notifications for DAO operations like Where/Data.
type Notifications struct {
	g.Meta     `orm:"table:notifications, do:true"`
	Id         any         // 提醒ID
	UserId     any         // 用户ID
	QuestionId any         // 问题ID
	ReplyToId  any         // 回复问题的ID
	AnswerId   any         // 问题ID
	Type       any         // 提醒类型（新提问、新回复、新回答）
	IsRead     any         // 是否已读
	CreatedAt  *gtime.Time // 创建时间
	DeletedAt  *gtime.Time // 删除时间
}
