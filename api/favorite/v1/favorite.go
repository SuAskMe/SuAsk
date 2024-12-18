package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"suask/internal/model"
)

type FavoriteReq struct {
	g.Meta `path:"/favorite/get" method:"GET" tags:"Favorite" summary:"查询收藏列表"`
}

type FavoriteRes struct {
	FavoriteQuestionList []model.FavoriteQuestion `json:"favorite_question_list"`
}

type PageFavoriteReq struct {
	g.Meta  `path:"/favorite/getPage" method:"GET" tags:"Favorite" summary:"分页查询收藏列表"`
	PageIdx int `json:"page_idx"`
}

type PageFavoriteRes struct {
	PageFavoriteQuestionList []model.FavoriteQuestion `json:"page_favorite_question_list"`
	Total                    int                      `json:"total" dc:"总问题数"`
	Size                     int                      `json:"size" dc:"每页问题数"`
	RemainPage               int                      `json:"remain_page" dc:"剩余页数"`
	PageNum                  int                      `json:"page_num" dc:"总页数"`
}

type DeleteFavoriteReq struct {
	g.Meta `path:"/favorite/delete" method:"DELETE" tags:"Favorite" summary:"删除收藏"`
	Id     int `json:"id"`
}

type DeleteFavoriteRes struct {
	String string `json:"string"`
}
