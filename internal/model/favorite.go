package model

import "github.com/gogf/gf/v2/os/gtime"

// model层：controller和logic之间的数据规范

// “获取收藏”相关结构体

type FavoriteQuestionInPut struct { // controller -> logic
	Id int `json:"id" v:"required" dc:"用户ID"`
}

type FavoriteQuestionOutPut struct { // logic -> controller
	FavoriteQuestionList []FavoriteQuestion `json:"Favorite_question" dc:"返回收藏问题列表"`
	//RemainPage int            `json:"remain_page" dc:"剩余页码数量"`
}

type FavoriteQuestion struct {
	ID         int         `json:"id" dc:"问题ID"`
	Title      string      `json:"title" dc:"标题"`
	Contents   string      `json:"contents" dc:"问题内容"`
	Views      int         `json:"views" dc:"浏览量"`
	FavoriteAt *gtime.Time `json:"favorite_at" dc:"收藏时间"`
}

// “删除收藏”相关结构体

type DeleteFavoriteInput struct {
	Id int `json:"id" v:"required" dc:"收藏Id"`
}

type DeleteFavoriteOutput struct {
	String string `json:"string" dc:"提示信息"`
}
