package model

import "github.com/gogf/gf/v2/os/gtime"

// model层：controller和logic之间的数据规范

// “获取收藏”相关结构体

type FavoriteQuestionInPut struct { // controller -> logic
	Id int `json:"id" v:"required" dc:"用户ID"`
}

type FavoriteQuestionOutPut struct { // logic -> controller
	FavoriteQuestionList []FavoriteQuestion `json:"favorite_question" dc:"返回收藏问题列表"`
}

type PageFavoriteQuestionInPut struct {
	Id      int `json:"id" v:"required" dc:"用户ID"`
	PageIdx int `json:"page_idx"`
}

type PageFavoriteQuestionOutPut struct {
	PageFavoriteQuestionList []FavoriteQuestion `json:"page_favorite_question" dc:"返回收藏问题列表"`
	Total                    int                `json:"total" dc:"总问题数"`
	Size                     int                `json:"size" dc:"每页问题数"`
	PageNum                  int                `json:"page_num" dc:"总页数"`
	RemainPage               int                `json:"remain_page" dc:"剩余页数"`
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
