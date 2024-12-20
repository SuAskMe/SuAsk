package custom

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type PublicQuestions struct {
	g.Meta    `orm:"table:questions"`
	Id        int         `json:"id"        orm:"id"          description:"问题ID"`  // 问题ID
	Title     string      `json:"title"     orm:"title"       description:"问题标题"`  // 问题标题
	Contents  string      `json:"contents"  orm:"contents"    description:"问题内容"`  // 问题内容
	CreatedAt *gtime.Time `json:"created_at" orm:"created_at"  description:"创建时间"` // 创建时间
	Views     int         `json:"views"     orm:"views"       description:"浏览量"`   // 浏览量
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

type Avatars struct {
	g.Meta       `orm:"table:users"`
	UserId       int `json:"userId" orm:"id"`
	AvatarFileId int `json:"avatarFileId" orm:"avatar_file_id"`
}

type AnswerImage struct {
	g.Meta   `orm:"table:attachments"`
	AnswerId int `json:"answer_id" orm:"answer_id" description:"回答ID"`  // 问题ID
	FileID   int `json:"file_id" orm:"file_id"      description:"文件ID"` // 文件ID
}
