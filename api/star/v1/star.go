package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"suask/internal/model"
)

type StarReq struct {
	g.Meta `path:"/star/get" method:"GET" tag:"GetStar" summary:"查询收藏列表"`
}

type StarRes struct {
	StarQuestionList model.StarQuestionOutPut `json:"star_question_list"`
}
