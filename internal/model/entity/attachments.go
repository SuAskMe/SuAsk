// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Attachments is the golang structure for table attachments.
type Attachments struct {
	Id         int    `json:"id"         orm:"id"          description:"附件ID"`          // 附件ID
	QuestionId int    `json:"questionId" orm:"question_id" description:"问题ID"`          // 问题ID
	AnswerId   int    `json:"answerId"   orm:"answer_id"   description:"回答ID"`          // 回答ID
	Type       string `json:"type"       orm:"type"        description:"附件类型（目前仅支持图片）"` // 附件类型（目前仅支持图片）
	FileId     int    `json:"fileId"     orm:"file_id"     description:"文件ID"`          // 文件ID
}
