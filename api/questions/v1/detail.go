package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetDetailReq struct {
	g.Meta     `path:"/question/detail" method:"get" tags:"question" summary:"获取问题详情" description:"获取问题详情"`
	QuestionID int `json:"question_id"`
}

type GetDetailRes struct {
	Question model.QuestionBase        `json:"question"`
	Answers  []model.AnswerWithDetails `json:"answer_list"`
	CanReply bool                      `json:"can_reply"`
}
