package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetPageBase struct {
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
	Page     int `v:"required|min:1" json:"page"`
	UserID   int `v:"required|min:1" json:"user_id"`
}

type GetPageReq struct {
	g.Meta `path:"/questions/public/get" method:"get" tags:"public question" summary:"获取公开问题列表" description:"获取公开问题列表"`
	GetPageBase
}

type GetPageRes struct {
	QuestionList []model.PublicQuestion `json:"question_list"`
	RemainPage   int                    `json:"remain_page"`
}

type GetSearchKeywordsReq struct {
	g.Meta   `path:"/questions/public/get/keywords" method:"get" tags:"public question" summary:"搜索公开问题" description:"搜索公开问题"`
	Keyword  string `v:"required|length:1,100" json:"keyword"`
	SortType int    `v:"required|min:0|max:3" json:"sort_type"`
}

type GetSearchKeywordsRes struct {
	Words []struct {
		Value string `json:"value"`
	} `json:"words"`
}

type GetPageByKeywordReq struct {
	g.Meta  `path:"/questions/public/search" method:"get" tags:"public question" summary:"根据关键字获取公开问题列表" description:"根据关键字获取公开问题列表"`
	Keyword string `v:"length:1,20" json:"keyword"`
	GetPageBase
}

type GetPageByKeywordRes struct {
	QuestionList []model.PublicQuestion `json:"question_list"`
	RemainPage   int                    `json:"remain_page"`
}

type FavoriteReq struct {
	g.Meta     `path:"/questions/public/favorite" method:"post" tags:"public question" summary:"收藏公开问题" description:"收藏公开问题"`
	QuestionID int `v:"required|min:1" json:"question_id"`
	UserID     int `v:"required|min:1" json:"user_id"`
}

type FavoriteRes struct {
	IsFavorited bool `json:"is_favorited"`
}

// type UpvoteReq struct {
// 	g.Meta     `path:"/questions/public/upvote" method:"post" tags:"public question" summary:"点赞公开问题" description:"点赞公开问题"`
// 	QuestionID int `v:"required|min:1" json:"question_id"`
// 	UserID     int `v:"required|min:1" json:"user_id"`
// }

// type UpvoteRes struct {
// 	Upvotes   int  `json:"upvotes"`
// 	IsSuccess bool `json:"is_success"`
// }
