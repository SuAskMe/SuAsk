package model

import (
	"suask/internal/model/custom"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type MyHistoryQuestion struct {
	Id        int         `json:"id" dc:"问题ID"`
	Contents  string      `json:"contents" dc:"问题内容"`
	CreatedAt *gtime.Time `json:"create_at" dc:"创建时间"`
	Views     int         `json:"views" dc:"浏览量"`
	ImageURLs []string    `json:"image_urls" dc:"图片的url"`
}

type GetHistoryInput struct {
	UserId int `json:"user_id" dc:"发起提问者的ID"`
	Page   int `json:"page" dc:"历史提问的第几页"`
}

type GetHistoryOutput struct {
	Question   []MyHistoryQuestion `json:"history_question" dc:"返回历史提问列表"`
	RemainPage int                 `json:"remain_page" dc:"剩余页码数量"`
}

// 使用custom中定义的struct实现存储多表查询的结构体
type MultiQueryQuestions struct {
	g.Meta    `orm:"table:questions"`
	Id        int             `json:"id"        orm:"id"          description:"问题ID"`
	Contents  string          `json:"contents"  orm:"contents"    description:"问题内容"`
	CreatedAt *gtime.Time     `json:"createdAt" orm:"created_at"  description:"创建时间"`
	Views     int             `json:"views"     orm:"views"       description:"浏览量"`
	Images    []*custom.Image `json:"images"  orm:"with:question_id=id" description:"1"`
}
