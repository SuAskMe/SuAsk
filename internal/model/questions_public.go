package model

type PublicQuestion struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"` // 目前不支持标题
	Content       string   `json:"contents"`
	Views         int      `json:"views"`
	CreatedAt     int64    `json:"created_at"`
	ImageURLs     []string `json:"image_urls"`
	IsFavorited   bool     `json:"is_favorited"`
	IsUpvoted     bool     `json:"is_upvoted"`
	AnswerNum     int      `json:"answer_num"`
	AnswerAvatars []string `json:"answer_avatars"`
}

type GetInput struct {
	SortType int    `json:"sort_type"`
	Page     int    `json:"page"`
	UserID   int    `json:"user_id"`
	Keyword  string `json:"keyword"`
}

type GetOutput struct {
	Questions  []PublicQuestion `json:"questions"`
	RemainPage int              `json:"remain_page"`
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

type FavoriteInput struct {
	QuestionID int `json:"question_id"`
	UserID     int `json:"user_id"`
}
