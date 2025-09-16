package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetPageBase struct {
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
	Page     int `v:"required|min:1" json:"page"`
}

type GetHistoryPageReq struct {
	g.Meta `path:"/history" method:"GET" tags:"History" summary:"获取收藏列表"`
	GetPageBase
}

type GetHistoryPageRes struct {
	QuestionList []model.HistoryQuestion `json:"question_list"`
	RemainPage   int                     `json:"remain_page"`
}

type GetHistorySearchKeywordsReq struct {
	g.Meta   `path:"/history/keywords" method:"GET" tags:"History" summary:"搜索收藏"`
	Keyword  string `v:"required|length:2,100" json:"keyword"`
	SortType int    `v:"required|min:0|max:3" json:"sort_type"`
}

type GetHistorySearchKeywordsRes struct {
	Words []struct {
		Value string `json:"value"`
	} `json:"words"`
}

type GetHistoryPageByKeywordReq struct {
	g.Meta  `path:"/history/search" method:"GET" tags:"History" summary:"根据关键字获取收藏列表"`
	Keyword string `v:"length:2,100" json:"keyword"`
	GetPageBase
}

type GetHistoryPageByKeywordRes struct {
	QuestionList []model.HistoryQuestion `json:"question_list"`
	RemainPage   int                     `json:"remain_page"`
}
