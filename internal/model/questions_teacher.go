package model

type TeacherQuestion struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"contents"`
	Views       int      `json:"views"`
	CreatedAt   int64    `json:"created_at"`
	ImageURLs   []string `json:"image_urls"`
	IsFavorited bool     `json:"is_favorited"`
}

type GetBaseOfTeacherInput struct {
	SortType  int    `json:"sort_type"`
	Page      int    `json:"page"`
	Keyword   string `json:"keyword"`
	TeacherID int    `json:"teacher_id"`
}

type GetBaseOfTeacherOutput struct {
	QuestionIDs []int             `json:"question_ids"`
	IdMap       map[int]int       `json:"id_map"`
	Questions   []TeacherQuestion `json:"questions"`
	RemainPage  int               `json:"remain_page"`
}

type GetKeywordsOfTeacherInput struct {
	Keyword   string `json:"keyword"`
	SortType  int    `json:"sort_type"`
	TeacherID int    `json:"teacher_id"`
}
