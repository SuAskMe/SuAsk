package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetPublicQuestionsReq struct {
	g.Meta   `path:"/questions/public/get" method:"get"`
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
	Page     int `v:"required|min:1" json:"page"`
	UserID   int `v:"required|min:1" json:"user_id"`
}

type GetPublicQuestionsRes struct {
	QuestionList []model.PublicQuestion `json:"question_list"`
	RemainPage   int                    `json:"remain_page"`
}
