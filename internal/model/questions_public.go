package model

import "github.com/gogf/gf/v2/encoding/gbinary"

type PublicQuestion struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"contents"`
	Views         int      `json:"views"`
	CreatedAt     int64    `json:"created_at"`
	ImageURLs     []string `json:"image_urls"`
	IsFavorite    bool     `json:"is_favorite"`
	AnswerNum     int      `json:"answer_num"`
	AnswerAvatars []string `json:"answer_avatars"`
}

type AddQuestionInput struct {
	SrcUserID interface{} `json:"src_user_id" orm:"src_user_id"`
	DstUserID interface{} `json:"dst_user_id" orm:"dst_user_id"`
	Title     string      `json:"title" orm:"title"`
	Content   string      `json:"content" orm:"content"`
	IsPrivate gbinary.Bit `json:"is_private" orm:"is_private"`
}

type AddQuestionOutput struct {
	ID int `json:"id"`
}

type GetBaseInput struct {
	SortType int    `json:"sort_type"`
	Page     int    `json:"page"`
	Keyword  string `json:"keyword"`
}

type GetBaseOutput struct {
	QuestionIDs []int            `json:"question_ids"`
	IdMap       map[int]int      `json:"id_map"`
	Questions   []PublicQuestion `json:"questions"`
	RemainPage  int              `json:"remain_page"`
}

type GetKeywordsInput struct {
	Keyword  string `json:"keyword"`
	SortType int    `json:"sort_type"`
}

type Keywords struct {
	Value string `json:"value" orm:"title"`
}

type GetKeywordsOutput struct {
	Words []Keywords `json:"words"`
}

type GetAnswersInput struct {
	QuestionIDs []int `json:"question_ids"`
}

type GetAnswersOutput struct {
	CountMap   map[int]int   `json:"count_map"`
	AvatarsMap map[int][]int `json:"avatars_map"`
}
