package model

import "github.com/gogf/gf/v2/os/gtime"

// model层：controller和logic之间的数据规范

// “获取收藏”相关结构体

type StarQuestionInPut struct { // controller -> logic
	Id int `json:"id" v:"required" dc:"用户ID"`
}

type StarQuestionOutPut struct { // logic -> controller
	StarQuestionList []StarQuestion `json:"star_question" dc:"返回收藏问题列表"`
	//RemainPage int            `json:"remain_page" dc:"剩余页码数量"`
}

type StarQuestion struct {
	ID       int         `json:"id" dc:"问题ID"`
	Title    string      `json:"title" dc:"标题"`
	Contents string      `json:"contents" dc:"问题内容"`
	Views    int         `json:"views" dc:"浏览量"`
	StarAt   *gtime.Time `json:"star_at" dc:"收藏时间"`
}

// “删除收藏”相关结构体

type DeleteStarInput struct {
	Id int `json:"id" v:"required" dc:"收藏Id"`
}

type DeleteStarOutput struct {
	String string `json:"string" dc:"提示信息"`
}
