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

type DeleteFavoriteReq struct {
	g.Meta `path:"/favorite/delete" method:"DELETE" tags:"Favorite" summary:"删除收藏"`
	Id     int `json:"id"`
}

type DeleteFavoriteRes struct {
	String string `json:"string"`
}
