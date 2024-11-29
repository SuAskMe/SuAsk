// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Notifications is the golang structure of table notifications for DAO operations like Where/Data.
type Notifications struct {
	g.Meta     `orm:"table:notifications, do:true"`
	Id         interface{} // 提醒ID
	UserId     interface{} // 用户ID
	QuestionId interface{} // 问题ID
	Type       interface{} // 提醒类型（新提问或新回复）
}
