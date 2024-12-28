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
	Id         interface{} // 提醒ID
	UserId     interface{} // 用户ID
	QuestionId interface{} // 问题ID
	AnswerId   interface{} // 问题ID
	Type       interface{} // 提醒类型（新提问、新回复、新回答）
	IsRead     interface{} // 是否已读
	CreatedAt  *gtime.Time // 创建时间
	DeletedAt  *gtime.Time // 删除时间
}
