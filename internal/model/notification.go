package model

type AddNotificationInput struct {
	UserId     int    `json:"user_id"     orm:"user_id"     dc:"用户ID"`
	QuestionId int    `json:"question_id" orm:"question_id" dc:"问题ID"`
	ReplyToId  int    `json:"reply_to_id"  orm:"reply_to_id" dc:"回复问题的ID"`
	AnswerId   int    `json:"answer_id"   orm:"answer_id"   dc:"回复ID"`
	Type       string `json:"type"       orm:"type"        dc:"提醒类型（新提问、新回复、新回答）"`
}

type NotificationQuestion struct {
	Id        int    `json:"id"        orm:"id"          description:"问题ID"`
	Title     string `json:"title"     orm:"title"       description:"问题标题"`
	Contents  string `json:"contents"  orm:"contents"    description:"问题内容"`
	DstUserId int    `json:"dstUserId" orm:"dst_user_id" description:"被提问的用户ID，为空时问大家，不为空时问教师"`
}

type NotificationAnswer struct {
	Id       int    `json:"id"         orm:"id"          description:"回答ID"`
	UserId   int    `json:"userId"     orm:"user_id"     description:"用户ID"`
	Contents string `json:"contents"   orm:"contents"    description:"回答内容"`
}

type NotificationUser struct {
	Id           int    `json:"id"           orm:"id"             description:"用户ID"`
	Nickname     string `json:"nickname"     orm:"nickname"       description:"昵称"`
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"`
}

type AddNotificationOutput struct {
	Id int `json:"id"         orm:"id"          dc:"提醒ID"`
}

type GetNotificationsInput struct {
	UserId int `json:"user_id"     orm:"user_id"     dc:"用户ID"`
}

type NotificationBase struct {
	Id              int64  `json:"id"         orm:"id"          dc:"提醒ID"`
	QuestionId      int    `json:"question_id" orm:"question_id" dc:"问题ID"`
	QuestionTitle   string `json:"question_title"                dc:"问题标题"`
	QuestionContent string `json:"question_content"              dc:"问题内容"`
	IsRead          bool   `json:"is_read"     orm:"is_read"     dc:"是否已读"`
	CreatedAt       int64  `json:"created_at"  orm:"created_at"  dc:"提醒创建时间"`
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

type UpdateAoQInput struct {
	UserID     int `json:"user_id"     orm:"user_id"     dc:"用户ID"`
	QuestionID int `json:"question_id" orm:"question_id" dc:"问题ID"`
}

type UpdateAoQOutput struct {
	QuestionID int  `json:"question_id" orm:"question_id" dc:"问题ID"`
	IsRead     bool `json:"isRead"     orm:"is_read"     dc:"是否已读"`
}

type DeleteNotificationInput struct {
	Id int `json:"id"         orm:"id"          dc:"提醒ID"`
}

type DeleteNotificationOutput struct{}

type NewNotificationCountInput struct {
	UserId int `json:"user_id" dc:"用户ID"`
}

type NewNotificationCountOutput struct {
	NewQuestionCount int `json:"new_question_count" dc:"新问题数目"`
	NewReplyCount    int `json:"new_reply_count" dc:"新回复数目"`
	NewAnswerCount   int `json:"new_answer_count" dc:"新回答数目"`
}
