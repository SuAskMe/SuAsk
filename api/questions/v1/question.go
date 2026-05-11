package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// GetPageBase 列表接口的公共分页参数，原来定义在已删除的 public.go 里，
// 现在放在这里给 teacher.go 继续复用。
type GetPageBase struct {
	SortType int `v:"required|min:0|max:3" json:"sort_type"`
	Page     int `v:"required|min:1" json:"page"`
}

type AddQuestionReq struct {
	g.Meta    `path:"/questions/add" method:"post" tags:"Question" summary:"添加一个问题"`
	DstUserId int                 `json:"dst_user_id" v:"required|min:1#请指定目标老师"`
	Title     string              `json:"title" v:"required"`
	Content   string              `json:"content" v:"required"`
	IsPrivate bool                `json:"is_private"`
	Files     []*ghttp.UploadFile `json:"files"`
}

type AddQuestionRes struct {
	Id int `json:"id"`
}
