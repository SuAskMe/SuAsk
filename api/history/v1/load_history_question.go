package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type LoadHistoryQuestionReq struct {
	g.Meta `path:"/history" method:"GET" tag:"history" summary:"获取历史提问请求"`
	UserId int `v:"required|min:1" json:"user_id" dc:"发起获取历史提问请求的UserId"`
	Page   int `v:"required|min:1" json:"page" dc:"历史提问的第几页"`
}

type LoadHistoryQuestionRes struct {
	HistoryQuestionList []model.MyHistoryQuestion `json:"question_list" dc:"后端返回的历史提问列表"`
	RemainPage          int                       `json:"remain_page" dc:"剩余页码数量"`
}
