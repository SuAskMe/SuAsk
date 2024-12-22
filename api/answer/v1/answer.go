package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"suask/internal/model"
)

type GetDetailReq struct {
	g.Meta     `path:"/answer" method:"get" tags:"Answer" summary:"获取问题回复" description:"获取问题回复"`
	QuestionID int `json:"question_id"`
}

type GetDetailRes struct {
	Question model.QuestionBase        `json:"question"`
	Answers  []model.AnswerWithDetails `json:"answer_list"`
	CanReply bool                      `json:"can_reply"`
}

type UpvoteReq struct {
	g.Meta     `path:"/answer/upvote" method:"post" tags:"Answer" summary:"点赞回复"`
	QuestionID int `v:"required|min:1" json:"question_id"`
	AnswerID   int `v:"required|min:1" json:"answer_id"`
}

type UpvoteRes struct {
	IsUpvoted bool `json:"is_upvoted"`
	UpvoteNum int  `json:"upvote_num"`
}

type AddAnswerReq struct {
	g.Meta     `path:"/answer/add" method:"post" tags:"Answer" summary:"添加一个回答"`
	QuestionId int                 `json:"question_id" v:"required"`
	Content    string              `json:"content" v:"required"`
	Files      []*ghttp.UploadFile `json:"files"`
}

type AddAnswerRes struct {
	Id int `json:"id"`
}
