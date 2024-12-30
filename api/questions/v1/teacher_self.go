package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type QFMBase struct {
	QFMList    []model.QFM `json:"question_list"`
	RemainPage int         `json:"remain_page"`
}

type GetQFMReq struct {
	g.Meta   `path:"/teacher/question/all" method:"GET" tags:"Teacher Self" summary:"获取对我的提问"`
	Page     int `json:"page" v:"required|min:1"`
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
}

type GetQFMRes struct {
	QFMBase
}

type GetQFMSearchKeywordsReq struct {
	// 只允许搜索所有的提问，不允许按分类搜索
	g.Meta   `path:"/teacher/question/keywords" method:"GET" tags:"Teacher Self" summary:"获取对我的提问的关键字"`
	Keyword  string `v:"required|length:1,100" json:"keyword"`
	SortType int    `v:"required|min:0|max:3" json:"sort_type"`
}

type GetQFMSearchKeywordsRes struct {
	Words []struct {
		Value string `json:"value"`
	} `json:"words"`
}

type SearchQFMReq struct {
	// 只允许搜索所有的提问，不允许按分类搜索
	g.Meta   `path:"/teacher/question/search" method:"GET" tags:"Teacher Self" summary:"搜索对我的提问"`
	Keyword  string `v:"required|length:1,100" json:"keyword"`
	SortType int    `v:"required|min:0|max:3" json:"sort_type"`
	Page     int    `json:"page" v:"required|min:1"`
}

type SearchQFMRes struct {
	QFMBase
}

type GetQFMUnansweredReq struct {
	g.Meta   `path:"/teacher/question/unanswered" method:"GET" tags:"Teacher Self" summary:"获取未回复提问"`
	Page     int `json:"page" v:"required|min:1"`
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
}

type GetQFMUnansweredRes struct {
	QFMBase
}

type GetQFMAnsweredReq struct {
	g.Meta   `path:"/teacher/question/answered" method:"GET" tags:"Teacher Self" summary:"获取已回复提问"`
	Page     int `json:"page" v:"required|min:1"`
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
}

type GetQFMAnsweredRes struct {
	QFMBase
}

type GetQFMTopReq struct {
	g.Meta `path:"/teacher/question/top" method:"GET" tags:"Teacher Self" summary:"获取置顶提问"`
}

type GetQFMTopRes struct {
	QFMBase
}

type PinQFMReq struct {
	g.Meta     `path:"/teacher/question/pin" method:"POST" tags:"Teacher Self" summary:"置顶提问"`
	QuestionId int `json:"question_id" v:"required|min:1"`
}

type PinQFMRes struct {
	IsPinned bool `json:"is_pinned"`
}
