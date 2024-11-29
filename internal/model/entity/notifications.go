// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Notifications is the golang structure for table notifications.
type Notifications struct {
	Id         int    `json:"id"         orm:"id"          description:"提醒ID"`          // 提醒ID
	UserId     int    `json:"userId"     orm:"user_id"     description:"用户ID"`          // 用户ID
	QuestionId int    `json:"questionId" orm:"question_id" description:"问题ID"`          // 问题ID
	Type       string `json:"type"       orm:"type"        description:"提醒类型（新提问或新回复）"` // 提醒类型（新提问或新回复）
}
