package model

type PublicQuestion struct {
	ID            int      `json:"id"`
	Title         string   `json:"title"` // 目前不支持标题
	Content       string   `json:"contents"`
	Views         int      `json:"views"`
	CreatedAt     int64    `json:"created_at"`
	ImageURLs     []string `json:"image_urls"`
	IsFavorited   bool     `json:"is_favorited"`
	AnswerNum     int      `json:"answer_num"`
	AnswerAvatars []string `json:"answer_avatars"`
}

type GetBaseInput struct {
	SortType int    `json:"sort_type"`
	Page     int    `json:"page"`
	UserID   int    `json:"user_id"`
	Keyword  string `json:"keyword"`
}

type GetBaseOutput struct {
	QuestionIDs []int            `json:"question_ids"`
	IdMap       map[int]int      `json:"id_map"`
	Questions   []PublicQuestion `json:"questions"`
	RemainPage  int              `json:"remain_page"`
}

type GetImagesInput struct {
	QuestionIDs []int `json:"question_ids"`
}

type GetImagesOutput struct {
	ImageMap map[int][]int `json:"image_map"`
}

type GetAnswersInput struct {
	QuestionIDs []int `json:"question_ids"`
}

type GetAnswersOutput struct {
	CountMap   map[int]int   `json:"count_map"`
	AvatarsMap map[int][]int `json:"avatars_map"`
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

type FavoriteOutput struct {
	IsFavorited bool `json:"is_favorited"`
}
