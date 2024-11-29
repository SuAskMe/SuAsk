// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Answers is the golang structure for table answers.
type Answers struct {
	Id         int         `json:"id"         orm:"id"          description:"回答ID"`        // 回答ID
	UserId     int         `json:"userId"     orm:"user_id"     description:"用户ID"`        // 用户ID
	QuestionId int         `json:"questionId" orm:"question_id" description:"问题ID"`        // 问题ID
	InReplyTo  int         `json:"inReplyTo"  orm:"in_reply_to" description:"回复的回答ID，可为空"` // 回复的回答ID，可为空
	Contents   string      `json:"contents"   orm:"contents"    description:"回答内容"`        // 回答内容
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:"创建时间"`        // 创建时间
	Views      int         `json:"views"      orm:"views"       description:"浏览量"`         // 浏览量
}
