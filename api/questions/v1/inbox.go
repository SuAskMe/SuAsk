package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// --- 统一收件箱（替代 /teacher/question/* 系列接口） ---
// 任何有 dst_user_id 指向自己的用户都能用（目前只有老师）。
// 前端迁移完后，旧的 /teacher/question/* 接口可以删除。

type InboxReq struct {
	g.Meta   `path:"/questions/inbox" method:"GET" tags:"Inbox" summary:"我的收件箱"`
	Page     int    `json:"page" v:"required|min:1"`
	SortType int    `json:"sort_type" v:"required|min:0|max:3"`
	Tag      string `json:"tag"` // all / answered / unanswered / pinned；空=all
}

type InboxRes struct {
	QuestionList []model.QFM `json:"question_list"`
	RemainPage   int         `json:"remain_page"`
}

type InboxKeywordsReq struct {
	g.Meta  `path:"/questions/inbox/keywords" method:"GET" tags:"Inbox" summary:"收件箱关键字"`
	Keyword string `json:"keyword" v:"required|length:2,100#请输入关键词|关键词长度需为 2-100 个字符"`
}

type InboxKeywordsRes struct {
	Words []struct {
		Value string `json:"value"`
	} `json:"words"`
}

type InboxSearchReq struct {
	g.Meta   `path:"/questions/inbox/search" method:"GET" tags:"Inbox" summary:"收件箱搜索"`
	Keyword  string `json:"keyword" v:"required|length:2,100#请输入关键词|关键词长度需为 2-100 个字符"`
	Page     int    `json:"page" v:"required|min:1"`
	SortType int    `json:"sort_type" v:"required|min:0|max:3"`
}

type InboxSearchRes struct {
	QuestionList []model.QFM `json:"question_list"`
	RemainPage   int         `json:"remain_page"`
}

type PinReq struct {
	g.Meta     `path:"/questions/pin" method:"POST" tags:"Inbox" summary:"置顶/取消置顶"`
	QuestionId int `json:"question_id" v:"required|min:1"`
}

type PinRes struct {
	IsPinned bool `json:"is_pinned"`
}
