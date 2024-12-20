package model

type QuestionBase struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"contents"`
	Views     int      `json:"views"`
	CreatedAt int64    `json:"created_at"`
	ImageURLs []string `json:"image_urls"`
}

type AnswerWithDetails struct {
	Id         int      `json:"id"`          // 回答ID
	UserId     int      `json:"user_id"`     // 用户ID
	UserAvatar string   `json:"user_avatar"` // 用户头像
	Contents   string   `json:"contents"`    // 回答内容
	CreatedAt  int64    `json:"created_at"`  // 创建时间
	Upvotes    int      `json:"upvotes"`     // 点赞量
	ImageURLs  []string `json:"image_urls"`
}

type GetQuestionBaseInput struct {
	QuestionId int `json:"question_id"`
}

type GetQuestionBaseOutput struct {
	Question *QuestionBase `json:"question"`
	CanReply bool
}

type GetAnswerDetailInput struct {
	QuestionId int `json:"question_id"`
}

type GetAnswerDetailOutput struct {
	Answers    []AnswerWithDetails `json:"answers"`
	IdMap      map[int]int         `json:"id_map"`
	AvatarsMap map[int][]int       `json:"avatars_map"`
	ImageMap   map[int][]int       `json:"image_map"`
}
