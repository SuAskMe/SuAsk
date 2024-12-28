package model

type AddNotificationInput struct {
	UserId     int         `json:"userId"     orm:"user_id"     dc:"用户ID"`
	QuestionId int         `json:"questionId" orm:"question_id" dc:"问题ID"`
	ReplyToId  interface{} `json:"replyToId"  orm:"reply_to_id" dc:"回复问题的ID"`
	AnswerId   interface{} `json:"answerId"   orm:"answer_id"   dc:"回复ID"`
	Type       string      `json:"type"       orm:"type"        dc:"提醒类型（新提问、新回复、新回答）"`
}

type AddNotificationOutput struct {
	Id int `json:"id"         orm:"id"          dc:"提醒ID"`
}

type GetNotificationsInput struct {
	UserId int `json:"userId"     orm:"user_id"     dc:"用户ID"`
}

type NotificationBase struct {
	Id              int    `json:"id"         orm:"id"          dc:"提醒ID"`
	QuestionId      int    `json:"questionId" orm:"question_id" dc:"问题ID"`
	QuestionTitle   string `json:"question_title"                dc:"问题标题"`
	QuestionContent string `json:"question_content"              dc:"问题内容"`
	IsRead          bool   `json:"isRead"     orm:"is_read"     dc:"是否已读"`
	CreatedAt       int64  `json:"createdAt"  orm:"created_at"  dc:"提醒创建时间"`
}

type NotificationNewQuestion struct {
	NotificationBase
	UserAvatar string `json:"user_avatar" dc:"提问者的头像 URL"`
	UserName   string `json:"user_name" dc:"提问者的昵称"`
	UserId     int    `json:"user_id" dc:"提问者的ID"`
}

type NotificationNewAnswer struct {
	NotificationBase
	AnswerId         int    `json:"answer_id"   orm:"answer_id"   dc:"回答ID"`
	AnswerContent    string `json:"answer_content"                dc:"回答内容"`
	RespondentAvatar string `json:"respondent_avatar" dc:"回复者的头像 URL"`
	RespondentName   string `json:"respondent_name" dc:"回复者的昵称"`
	RespondentId     int    `json:"respondent_id" dc:"回复者的Id"`
}

type NotificationNewReply struct {
	NotificationBase
	ReplyToId        int    `json:"answer_id"   orm:"answer_id"   dc:"本人的回答ID"`
	ReplyToContent   string `json:"answer_content"                dc:"本人的回答内容"`
	RespondentAvatar string `json:"respondent_avatar" dc:"回复者的头像 URL"`
	RespondentName   string `json:"respondent_name" dc:"回复者的昵称"`
	RespondentId     int    `json:"respondent_id" dc:"回复者的Id"`
	AnswerId         int    `json:"reply_id" dc:"回复ID"`
	AnswerContent    string `json:"reply_content" dc:"回复的内容"`
}

type GetNotificationsOutput struct {
	NewQuestion []NotificationNewQuestion `json:"new_question"`
	NewReply    []NotificationNewReply    `json:"new_reply"`
	NewAnswer   []NotificationNewAnswer   `json:"new_answer"`
}

type UpdateNotificationInput struct {
	Id int `json:"id"         orm:"id"          dc:"提醒ID"`
}

type UpdateNotificationOutput struct {
	Id     int  `json:"id"         orm:"id"          dc:"提醒ID"`
	IsRead bool `json:"isRead"     orm:"is_read"     dc:"是否已读"`
}

type DeleteNotificationInput struct {
	Id int `json:"id"         orm:"id"          dc:"提醒ID"`
}

type DeleteNotificationOutput struct{}
