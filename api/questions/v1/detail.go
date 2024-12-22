package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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

type UpvoteReq struct {
	g.Meta     `path:"/questions/public/favorite" method:"post" tags:"public question" summary:"收藏公开问题" description:"收藏公开问题"`
	QuestionID int `v:"required|min:1" json:"question_id"`
	AnswerID   int `v:"required|min:1" json:"answer_id"`
}

type UpvoteRes struct {
	IsUpvoted bool `json:"is_upvoted"`
	UpvoteNum int  `json:"upvote_num"`
}

type AddAnswerReq struct {
	g.Meta     `path:"/question/detail/add" method:"post" tags:"question" summary:"添加一个问题"`
	QuestionId int                 `json:"question_id"`
	Content    string              `json:"content" v:"required"`
	Files      []*ghttp.UploadFile `json:"files"`
}

type AddAnswerRes struct {
	Success bool `json:"success"`
}
