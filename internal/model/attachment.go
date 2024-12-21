package model

type AddAttachmentInput struct {
	QuestionId interface{} `json:"question_id" orm:"question_id"`
	AnswerId   interface{} `json:"answer_id" orm:"answer_id"`
	Type       string      `json:"type" orm:"type"`
	FileId     []int       `json:"file_id_list"`
}
type AddAttachmentOutput struct {
	Id []int `json:"id_list" orm:"id"`
}
