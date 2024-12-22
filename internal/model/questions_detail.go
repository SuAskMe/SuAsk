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
	Id          int      `json:"id"`          // 回答ID
	UserId      int      `json:"user_id"`     // 用户ID
	UserAvatar  string   `json:"user_avatar"` // 用户头像
	Contents    string   `json:"contents"`    // 回答内容
	CreatedAt   int64    `json:"created_at"`  // 创建时间
	Upvotes     int      `json:"upvotes"`     // 点赞量
	ImageURLs   []string `json:"image_urls"`
	IsUpvoted   bool     `json:"is_upvoted"`
	TeacherName string   `json:"teacher_name"`
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

type AddViewInput struct {
	QuestionId int `json:"question_id"`
}

type AddViewOutput struct {
	Views int `json:"views"`
}

type UpvoteInput struct {
	QuestionId int `json:"question_id"`
	AnswerId   int `json:"answer_id"`
}

type UpvoteOutput struct {
	IsUpvoted bool `json:"is_upvoted"`
	Upvotes   int  `json:"upvote_num"`
}

type AddAnswerInput struct {
	QuestionId int    `json:"question_id" orm:"dst_user_id"`
	Content    string `json:"content" orm:"content"`
}

type AddAnswerOutput struct {
	Id int `json:"id"`
}
