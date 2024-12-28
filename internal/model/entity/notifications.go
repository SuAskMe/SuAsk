// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Notifications is the golang structure for table notifications.
type Notifications struct {
	Id         int         `json:"id"         orm:"id"          description:"提醒ID"`              // 提醒ID
	UserId     int         `json:"userId"     orm:"user_id"     description:"用户ID"`              // 用户ID
	QuestionId int         `json:"questionId" orm:"question_id" description:"问题ID"`              // 问题ID
	AnswerId   int         `json:"answerId"   orm:"answer_id"   description:"问题ID"`              // 问题ID
	Type       string      `json:"type"       orm:"type"        description:"提醒类型（新提问、新回复、新回答）"` // 提醒类型（新提问、新回复、新回答）
	IsRead     bool        `json:"isRead"     orm:"is_read"     description:"是否已读"`              // 是否已读
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`              // 创建时间
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:"删除时间"`              // 删除时间
}
