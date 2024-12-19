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

	// Images []*Image `json:"images" orm:"with:question_id=id" description:"图片"` // 图片
}

type Image struct {
	g.Meta     `orm:"table:attachments"`
	QuestionId int `json:"question_id" orm:"question_id" description:"问题ID"` // 问题ID
	FileID     int `json:"file_id" orm:"file_id"      description:"文件ID"`    // 文件ID
}

type QuestionAnswers struct {
	g.Meta `orm:"table:answers"`
	Id     int          `orm:"id"          description:"回答ID"`    // 回答ID
	Users  *AnswerUsers `orm:"with:user_id=id" description:"回答者"` // 回答者

}

type AnswerUsers struct {
	g.Meta `orm:"table:users"`
	Id     int `orm:"id"         `
	Avatar int `orm:"avatar_file_id" description:"回答者头像文件ID"` // 回答者头像文件ID
}

type MyFavorites struct {
	g.Meta     `orm:"table:favorites"`
	QuestionId int `json:"question_id" orm:"question_id" description:"问题ID"` // 问题ID
}
