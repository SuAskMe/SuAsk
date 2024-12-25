package model

import "github.com/gogf/gf/v2/os/gtime"

type Favorite struct {
	Id         int         `json:"id"          description:"收藏（置顶）ID"` // 收藏（置顶）ID
	UserId     int         `json:"userId"      description:"用户ID"`     // 用户ID
	QuestionId int         `json:"questionId"  description:"问题ID"`     // 问题ID
	CreatedAt  *gtime.Time `json:"createdAt"   description:"创建时间"`     // 创建时间
	Package    string      `json:"package"     description:"收藏夹"`      // 收藏夹
}

type GetFavoriteBaseInput struct {
	SortType int    `json:"sort_type"`
	Page     int    `json:"page"`
	Keyword  string `json:"keyword"`
}

type GetFavoriteBaseOutput struct {
	QuestionIDs []int            `json:"question_ids"`
	IdMap       map[int]int      `json:"id_map"`
	Questions   []PublicQuestion `json:"questions"`
	RemainPage  int              `json:"remain_page"`
}

type GetFavoriteKeywordsInput struct {
	Keyword  string `json:"keyword"`
	SortType int    `json:"sort_type"`
}

type FavoriteKeywords struct {
	Value string `json:"value" orm:"title"`
}

type GetFavoriteKeywordsOutput struct {
	Words []Keywords `json:"words"`
}

type GetFavoriteAnswersInput struct {
	QuestionIDs []int `json:"question_ids"`
}

type GetFavoriteAnswersOutput struct {
	CountMap   map[int]int   `json:"count_map"`
	AvatarsMap map[int][]int `json:"avatars_map"`
}
