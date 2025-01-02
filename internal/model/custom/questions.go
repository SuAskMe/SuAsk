package custom

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type Questions struct {
	g.Meta    `orm:"table:questions"`
	Id        int         `json:"id"        orm:"id"          description:"问题ID"`          // 问题ID
	Title     string      `json:"title"     orm:"title"       description:"问题标题"`          // 问题标题
	DstUserId int         `json:"dest_user_id"   orm:"dst_user_id"     description:"目标ID"` // 目标ID
	Contents  string      `json:"contents"  orm:"contents"    description:"问题内容"`          // 问题内容
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at"  description:"创建时间"`         // 创建时间
	Views     int         `json:"views"     orm:"views"       description:"浏览量"`           // 浏览量
	ReplyCnt  int         `json:"reply_cnt" orm:"reply_cnt"   description:"回复数量"`          // 回复数量
}

type Image struct {
	g.Meta     `orm:"table:attachments"`
	QuestionId int `json:"question_id" orm:"question_id" description:"问题ID"` // 问题ID
	FileID     int `json:"file_id" orm:"file_id"      description:"文件ID"`    // 文件ID
}

type MyFavorites struct {
	g.Meta     `orm:"table:favorites"`
	QuestionId int `json:"question_id" orm:"question_id" description:"问题ID"` // 问题ID
}

type MyUpvotes struct {
	g.Meta   `orm:"table:upvotes"`
	AnswerId int `json:"answer_id" orm:"answer_id" description:"回答ID"` // 回答ID
}

type UserInfo struct {
	g.Meta       `orm:"table:users"`
	UserId       int    `json:"userId" orm:"id"`
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id"`
	Role         string `json:"role" orm:"role"`
	Name         string `json:"name" orm:"name"`
	NickName     string `json:"nickName" orm:"nickname"`
}

type AnswerImage struct {
	g.Meta   `orm:"table:attachments"`
	AnswerId int `json:"answer_id" orm:"answer_id" description:"回答ID"`  // 问题ID
	FileID   int `json:"file_id" orm:"file_id"      description:"文件ID"` // 文件ID
}
