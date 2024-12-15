package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"suask/internal/model"
)

type StarReq struct {
	g.Meta `path:"/star/get" method:"GET" tags:"Star" summary:"查询收藏列表"`
}

type StarRes struct {
	StarQuestionList []model.StarQuestion `json:"star_question_list"`
}

type DeleteStarReq struct {
	g.Meta `path:"/star/delete" method:"DELETE" tags:"Star" summary:"删除收藏"`
	Id     int `json:"id"`
}

type DeleteStarRes struct {
	String string `json:"string"`
}
