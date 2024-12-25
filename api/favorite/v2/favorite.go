package v2

import (
	"github.com/gogf/gf/v2/frame/g"
	"suask/internal/model"
)

type GetPageBase struct {
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
	Page     int `v:"required|min:1" json:"page"`
}

type GetPageReq struct {
	g.Meta `path:"/favorites" method:"get" tags:"Favorite" summary:"获取公开问题列表" description:"获取公开问题列表"`
	GetPageBase
}

type GetPageRes struct {
	QuestionList []model.PublicQuestion `json:"favorite_list"`
	RemainPage   int                    `json:"remain_page"`
}

type GetSearchKeywordsReq struct {
	g.Meta   `path:"/favorites/keywords" method:"get" tags:"Favorite" summary:"搜索公开问题" description:"搜索公开问题"`
	Keyword  string `v:"required|length:1,100" json:"keyword"`
	SortType int    `v:"required|min:0|max:3" json:"sort_type"`
}

type GetSearchKeywordsRes struct {
	Words []struct {
		Value string `json:"value"`
	} `json:"words"`
}

type GetPageByKeywordReq struct {
	g.Meta  `path:"/favorites/search" method:"get" tags:"Favorite" summary:"根据关键字获取公开问题列表" description:"根据关键字获取公开问题列表"`
	Keyword string `v:"length:1,20" json:"keyword"`
	GetPageBase
}

type GetPageByKeywordRes struct {
	QuestionList []model.Favorite `json:"favorite_list"`
	RemainPage   int              `json:"remain_page"`
}
