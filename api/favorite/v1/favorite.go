package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"suask/internal/model"
)

type GetPageBase struct {
	SortType int `v:"required|min:0|max:1" json:"sort_type"`
	Page     int `v:"required|min:1" json:"page"`
}

type GetFavoritePageReq struct {
	g.Meta `path:"/favorites" method:"GET" tags:"Favorite" summary:"获取收藏列表"`
	GetPageBase
}

type GetFavoritePageRes struct {
	QuestionList []model.FavoriteQuestion `json:"favorite_list"`
	RemainPage   int                      `json:"remain_page"`
}

// type GetFavoriteSearchKeywordsReq struct {
// 	g.Meta   `path:"/favorites/keywords" method:"GET" tags:"Favorite" summary:"搜索收藏"`
// 	Keyword  string `v:"required|length:1,100" json:"keyword"`
// 	SortType int    `v:"required|min:0|max:1" json:"sort_type"`
// }

// type GetFavoriteSearchKeywordsRes struct {
// 	Words []struct {
// 		Value string `json:"value"`
// 	} `json:"words"`
// }

// type GetFavoritePageByKeywordReq struct {
// 	g.Meta  `path:"/favorites/search" method:"GET" tags:"Favorite" summary:"根据关键字获取收藏列表"`
// 	Keyword string `v:"length:1,20" json:"keyword"`
// 	GetPageBase
// }

// type GetFavoritePageByKeywordRes struct {
// 	QuestionList []model.FavoriteQuestion `json:"favorite_list"`
// 	RemainPage   int                      `json:"remain_page"`
// }

type FavoriteReq struct {
	g.Meta     `path:"/favorites" method:"POST" tags:"Favorite" summary:"收藏问题"`
	QuestionID int `v:"required|min:1" json:"question_id"`
}

type FavoriteRes struct {
	IsFavorite bool `json:"is_favorite"`
}
