package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type NotificationGetReq struct {
	g.Meta `path:"/notification" method:"GET" tags:"Notification" summary:"通过用户Id拿到通知"`
	UserId int `json:"user_id" dc:"用户ID"`
}

type NotificationGetRes struct {
	NewQuestion []model.NotificationNewQuestion `json:"new_question" dc:"新的问题"`
	NewReply    []model.NotificationNewReply    `json:"new_reply" dc:"新的回复"`
	NewAnswer   []model.NotificationNewAnswer   `json:"new_answer" dc:"新的回答"`
}

type NotificationUpdateReq struct {
	g.Meta `path:"/notification" method:"PUT" tags:"Notification" summary:"更新已读信息"`
	Id     int `json:"id" dc:"提醒ID"`
}

type NotificationUpdateRes struct {
	Id     int  `json:"id" dc:"提醒ID"`
	IsRead bool `json:"is_read" dc:"是否已读"`
}

type NotificationDeleteReq struct {
	g.Meta `path:"/notification" method:"DELETE" tags:"Notification" summary:"删除提醒"`
	Id     int `json:"id" dc:"提醒ID"`
}

type NotificationDeleteRes struct{}

type NotificationGetCountReq struct {
	g.Meta `path:"/notification/count" method:"GET" tags:"Notification" summary:"获取提醒数目"`
	UserId int `json:"user_id" dc:"用户ID"`
}
type NotificationGetCountRes struct {
	NewQuestionCount int `json:"new_question_count" dc:"新问题数目"`
	NewReplyCount    int `json:"new_reply_count" dc:"新回复数目"`
	NewAnswerCount   int `json:"new_answer_count" dc:"新回答数目"`
}
