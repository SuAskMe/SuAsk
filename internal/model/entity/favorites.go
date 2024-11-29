// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Favorites is the golang structure for table favorites.
type Favorites struct {
	Id         int `json:"id"         orm:"id"          description:"收藏（置顶）ID"` // 收藏（置顶）ID
	UserId     int `json:"userId"     orm:"user_id"     description:"用户ID"`     // 用户ID
	QuestionId int `json:"questionId" orm:"question_id" description:"问题ID"`     // 问题ID
	Ordinal    int `json:"ordinal"    orm:"ordinal"     description:"收藏序号"`     // 收藏序号
}
