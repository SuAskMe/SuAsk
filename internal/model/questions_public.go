package model

import "github.com/gogf/gf/os/gtime"

type PublicQuestion struct {
	ID            int         `json:"id"`
	Title         string      `json:"title"` // 目前不支持标题
	Content       *string     `json:"content"`
	Views         int         `json:"views"`
	CreatedAt     *gtime.Time `json:"created_at"`
	ImageURLs     []string    `json:"image_urls"`
	IsFavorited   bool        `json:"is_favourite"`
	IsUpvoted     bool        `json:"is_upvoted"`
	AnswerNum     int         `json:"answer_num"`
	AnswerAvatars []string    `json:"answer_avatars"`
}

type GetPublicQuestionInput struct {
	SortType int `json:"sort_type"`
	Page     int `json:"page"`
	UserID   int `json:"user_id"`
}

type GetPublicQuestionOutput struct {
	Questions  []PublicQuestion `json:"questions"`
	RemainPage int              `json:"remain_page"`
}
