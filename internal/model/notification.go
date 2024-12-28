package model

type AddNotificationInput struct {
	UserId     int         `json:"userId"     orm:"user_id"     description:"用户ID"`              // 用户ID
	QuestionId interface{} `json:"questionId" orm:"question_id" description:"问题ID"`              // 问题ID
	AnswerId   interface{} `json:"answerId"   orm:"answer_id"   description:"回复ID"`              // 回复ID
	Type       string      `json:"type"       orm:"type"        description:"提醒类型（新提问、新回复、新回答）"` // 提醒类型（新提问、新回复、新回答）
}

type AddNotificationOutput struct {
	Id int `json:"id"         orm:"id"          description:"提醒ID"`
}

type GetNotificationsInput struct {
	UserId int `json:"userId"     orm:"user_id"     description:"用户ID"` // 用户ID
}

type Notification struct {
	Id              int    `json:"id"         orm:"id"          description:"提醒ID"`
	QuestionId      int    `json:"questionId" orm:"question_id" description:"问题ID"`
	QuestionTitle   string `json:"questionTitle"                description:"问题标题"`
	QuestionContent string `json:"questionContent"              description:"问题内容"`
	AnswerId        int    `json:"answerId"   orm:"answer_id"   description:"回复ID"`
	AnswerContent   string `json:"answerContent"                description:"回复内容"`
	IsRead          bool   `json:"isRead"     orm:"is_read"     description:"是否已读"` // 是否已读
	CreatedAt       int64  `json:"createdAt"  orm:"created_at"  description:""`     //
}

type GetNotificationsOutput struct {
	NewQuestion []Notification `json:"new_question"`
	NewReply    []Notification `json:"new_reply"`
	NewAnswer   []Notification `json:"new_answer"`
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
