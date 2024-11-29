// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Upvotes is the golang structure for table upvotes.
type Upvotes struct {
	Id         int `json:"id"         orm:"id"          description:"点赞ID"` // 点赞ID
	UserId     int `json:"userId"     orm:"user_id"     description:"用户ID"` // 用户ID
	QuestionId int `json:"questionId" orm:"question_id" description:"问题ID"` // 问题ID
	AnswerId   int `json:"answerId"   orm:"answer_id"   description:"回复ID"` // 回复ID
}
