package v1

import (
	"suask/internal/model/custom"

	"github.com/gogf/gf/v2/frame/g"
)

type LoadHistoryQuestionReq struct {
	g.Meta `path:"/history" method:"GET" tag:"history" summary:"获取历史提问请求"`
	//UserId int `v:"required|min:1" json:"user_id" dc:"发起获取历史提问请求的UserId"`
}

type LoadHistoryQuestionRes struct {
	HistoryQuestionList []custom.PublicQuestions `json:"question_list" dc:"后端返回的历史提问列表"`
}

