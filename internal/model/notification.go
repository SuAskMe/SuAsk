package model

import "github.com/gogf/gf/v2/os/gtime"

type AddNotificationInput struct {
	UserId     int    `json:"userId"     orm:"user_id"     description:"用户ID"`          // 用户ID
	QuestionId int    `json:"questionId" orm:"question_id" description:"问题ID"`          // 问题ID
	AnswerId   int    `json:"answerId"   orm:"answer_id"   description:"问题ID"`          // 问题ID
	Type       string `json:"type"       orm:"type"        description:"提醒类型（新提问或新回复）"` // 提醒类型（新提问或新回复）
}

type AddNotificationOutput struct {
	Id int `json:"id"         orm:"id"          description:"提醒ID"`
}

type GetNotificationsInput struct {
	UserId int `json:"userId"     orm:"user_id"     description:"用户ID"` // 用户ID
}

type Notification struct {
	Id         int         `json:"id"         orm:"id"          description:"提醒ID"`          // 提醒ID
	QuestionId int         `json:"questionId" orm:"question_id" description:"问题ID"`          // 问题ID
	AnswerId   int         `json:"answerId"   orm:"answer_id"   description:"问题ID"`          // 问题ID
	Type       string      `json:"type"       orm:"type"        description:"提醒类型（新提问或新回复）"` // 提醒类型（新提问或新回复）
	IsRead     bool        `json:"isRead"     orm:"is_read"     description:"是否已读"`          // 是否已读
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`              //
}

type GetNotificationsOutput struct {
	Notifications []Notification `json:"notifications"`
}

type UpdateNotificationInput struct {
	Id int `json:"id"         orm:"id"          description:"提醒ID"` // 提醒ID
}

type UpdateNotificationOutput struct {
	Id     int  `json:"id"         orm:"id"          description:"提醒ID"` // 提醒ID
	IsRead bool `json:"isRead"     orm:"is_read"     description:"是否已读"` // 是否已读
}

type DeleteNotificationInput struct {
	Id int `json:"id"         orm:"id"          description:"提醒ID"` // 提醒ID
}

type DeleteNotificationOutput struct{}
