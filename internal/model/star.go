package model

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"suask/internal/model/do"
)

// model层：controller和logic之间的数据规范

// “获取收藏”相关结构体

type QuestionInfo struct {
	g.Meta   `orm:"table:questions"`
	ID       int    `json:"id" dc:"问题ID"`
	Title    string `json:"title"` // 目前不支持标题
	Contents string `json:"contents" dc:"问题内容"`
	Views    int    `json:"views" dc:"浏览量"`
}

type StarRelation struct {
	QuestionId int         `json:"question_id"`
	StarAt     *gtime.Time `json:"star_at"`

	Questions *do.Questions `orm:"with:id=question_id" json:"questions"`
}

type StarQuestionInPut struct { // controller -> logic
	Id int `json:"id" v:"required" dc:"用户ID"`
}

type StarQuestionOutPut struct { // logic -> controller
	Question []StarRelation `json:"star_question_list" dc:"返回收藏问题列表"`
	//RemainPage int            `json:"remain_page" dc:"剩余页码数量"`
}

// “删除收藏”相关结构体

type DeleteStarInput struct {
	Id int `json:"id" v:"required" dc:"收藏Id"`
}

type DeleteStarOutput struct {
	String string `json:"string" dc:"提示信息"`
}
