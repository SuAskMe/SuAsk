package model

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type StarQuestionInPut struct {
	r *ghttp.Request
}

type StarQuestionOutPut struct {
	ID        int         `json:"id"`
	Title     string      `json:"title"` // 目前不支持标题
	Content   *string     `json:"content"`
	Views     int         `json:"views"`
	CreatedAt *gtime.Time `json:"created_at"`
	ImageURLs []string    `json:"image_urls"`
	// 默认全部 IsFavorite = true
	AnswerNum     int      `json:"answer_num"`
	AnswerAvatars []string `json:"answer_avatars"`
}
