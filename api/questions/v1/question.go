package v1

import (
	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type AddQuestionReq struct {
	g.Meta `path:"/questions/add" method:"post" tags:"Question" summary:"添加一个问题"`
	//SrcUserId int                 `json:"src_user_id"`
	DstUserId int                 `json:"dst_user_id"`
	Title     string              `json:"title" v:"required"`
	Content   string              `json:"content" v:"required"`
	IsPrivate gbinary.Bit         `json:"is_private" v:"required"`
	Files     []*ghttp.UploadFile `json:"files"`
}

type AddQuestionRes struct {
	Id int `json:"id"`
}
