// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Viewed is the golang structure for table viewed.
type Viewed struct {
	UserId     int `json:"userId"     orm:"user_id"     description:"用户id"` // 用户id
	QuestionId int `json:"questionId" orm:"question_id" description:"问题id"` // 问题id
}
