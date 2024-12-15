package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"suask/internal/model"
)

type StarReq struct {
	g.Meta `path:"/star" method:"GET" tag:"Star" summary:"查询收藏列表"`
}

type StarRes struct {
	StarQuestionList []model.StarRelation `json:"star_question_list"`
}
