package model

type GetImagesInput struct {
	QuestionIDs []int `json:"question_ids"`
	AnswerIDs   []int `json:"answer_ids"`
}

type GetImagesOutput struct {
	ImageMap map[int][]int `json:"image_map"`
}

type FavoriteInput struct {
	QuestionID int `json:"question_id"`
}

type FavoriteOutput struct {
	IsFavorited bool `json:"is_favorited"`
}
