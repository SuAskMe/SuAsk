package custom

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/v2/frame/g"
)

type PublicQuestions struct {
	g.Meta    `orm:"table:questions"`
	Id        int         `json:"id"        orm:"id"          description:"问题ID"` // 问题ID
	Contents  string      `json:"contents"  orm:"contents"    description:"问题内容"` // 问题内容
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"  description:"创建时间"` // 创建时间
	Views     int         `json:"views"     orm:"views"       description:"浏览量"`  // 浏览量
	Upvotes   int         `json:"upvotes"   orm:"upvotes"     description:"点赞量"`  // 点赞量

	Images      []*Image          `json:"images"  orm:"with:question_id=id"`      // 图片附件
	IsUpvoted   []*UserUpvotes    `json:"is_upvoted" orm:"with:question_id=id"`   // 是否点赞
	IsFavorited []*UserFavorites  `json:"is_favorited" orm:"with:question_id=id"` // 是否收藏
	Answers     []*QuestionAnwers `json:"answers"   orm:"where:question_id=id"`   // 回答
}

type Image struct {
	g.Meta `orm:"table:attachments"`
	Type   string `json:"type"   orm:"type"         description:"附件类型（目前仅支持图片）"`
	FileID int    `json:"fileId" orm:"file_id"      description:"文件ID"` // 文件ID
}

type FileHash struct {
	g.Meta `orm:"table:files"`
	Name   string `json:"name"       orm:"name"        description:"文件名，不得包含非法字符例如斜杠"`  // 文件名，不得包含非法字符例如斜杠
	Hash   []byte `json:"hash"       orm:"hash"        description:"文件哈希，算法暂定为BLAKE2b"` // 文件哈希，算法暂定为BLAKE2b
}

type UserUpvotes struct {
	g.Meta `orm:"table:upvotes"`
	UserID int `json:"userId"     orm:"user_id"     description:"用户ID"` // 用户ID
}

type UserFavorites struct {
	g.Meta `orm:"table:favorites"`
	UserID int `json:"userId"     orm:"user_id"     description:"用户ID"` // 用户ID
}

type QuestionAnwers struct {
	g.Meta `orm:"table:answers"`
	Id     int          `json:"id"        orm:"id"          description:"回答ID"`    // 回答ID
	Users  *AnswerUsers `json:"users"     orm:"with:user_id=id" description:"回答者"` // 回答者
}

type AnswerUsers struct {
	g.Meta `orm:"table:users"`
	Avatar int `json:"avatar" orm:"avatar_file_id" description:"回答者头像文件ID"` // 回答者头像文件ID
}
