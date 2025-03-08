package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type LoadHistoryQuestionReq struct {
	g.Meta `path:"/history" method:"GET" tags:"History" summary:"获取历史提问请求"`
	Page   int `v:"required|min:1" json:"page" dc:"历史提问的第几页"`
}

type LoadHistoryQuestionRes struct {
	HistoryQuestionList []model.MyHistoryQuestion `json:"question_list" dc:"后端返回的历史提问列表"`
	Total               int                       `json:"total" dc:"总问题数"`
	Size                int                       `json:"size" dc:"每页问题数"`
	RemainPage          int                       `json:"remain_page" dc:"剩余页数"`
	PageNum             int                       `json:"page_num" dc:"总页数"`
}
