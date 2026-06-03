package model

// AddQuestionInput / AddQuestionOutput —— 之前在 questions_public.go 里，
// 随"问大家"模块下线一起挪到这里，作为统一的"新增问题"入参/出参。
type AddQuestionInput struct {
	SrcUserID int    `json:"src_user_id" orm:"src_user_id"`
	DstUserID int    `json:"dst_user_id" orm:"dst_user_id"`
	Title     string `json:"title" orm:"title"`
	Content   string `json:"content" orm:"content"`
}

type AddQuestionOutput struct {
	ID int `json:"id"`
}

// GetAnswersInput / GetAnswersOutput —— 批量拿"问题 → 回答者头像"的映射，
// 供收藏、历史、老师问题列表等多种列表页使用。
type GetAnswersInput struct {
	QuestionIDs []int `json:"question_ids"`
}

type GetAnswersOutput struct {
	AvatarsMap map[int][]int `json:"avatars_map"`
}

// Keyword 原来定义在 questions_public.go 里，被 history / teacher 等关键字搜索接口复用。
type Keyword struct {
	Value string `json:"value" orm:"title"`
}

// GetKeywordsOutput 关键词搜索的公共返回结构，history / teacher 接口都在用。
type GetKeywordsOutput struct {
	Words []Keyword `json:"words"`
}

// GetImagesInput / GetImagesOutput —— 批量拿问题附件图片 ID，列表接口复用
type GetImagesInput struct {
	QuestionIDs []int `json:"question_ids"`
}

type GetImagesOutput struct {
	ImageMap map[int][]int `json:"image_map"`
}

type GetQuestionListAssetsInput struct {
	QuestionIDs  []int       `json:"question_ids"`
	DstUserIDMap map[int]int `json:"dst_user_id_map"`
}

type GetQuestionListAssetsOutput struct {
	ImageURLMap     map[int][]string `json:"image_url_map"`
	AnswerAvatarMap map[int][]string `json:"answer_avatar_map"`
}

// FavoriteInput / FavoriteOutput —— 收藏/取消收藏
type FavoriteInput struct {
	QuestionID int `json:"question_id"`
}

type FavoriteOutput struct {
	IsFavorite bool `json:"is_favorite"`
}
