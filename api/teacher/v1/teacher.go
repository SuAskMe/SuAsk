package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type TeacherBase struct {
	Id           int    `json:"id"           orm:"id"           description:"教师Id"`   //
	Responses    int    `json:"responses"    orm:"responses"    description:"回复数"`    // 回复数
	Name         string `json:"name"         orm:"name"         description:"老师名字"`   // 老师名字
	AvatarUrl    string `json:"avatarUrl"    orm:"avatar_url"   description:"老师头像链接"` // 老师头像链接
	Introduction string `json:"introduction" orm:"introduction" description:"老师简介"`   // 老师简介
	Email        string `json:"email"        orm:"email"        description:"老师邮箱"`   // 老师邮箱
	Perm         string `json:"perm"         orm:"perm"         description:"提问箱权限"`  // 提问箱权限
}

type TeacherReq struct {
	g.Meta `path:"/info/teacher" method:"GET" tags:"Info" summary:"请求教师信息"`
}

type TeacherRes struct {
	TeacherList []TeacherBase `json:"teachers"`
}

type TeacherPinReq struct {
	g.Meta    `path:"/info/teacher/pin" method:"GET" tags:"Info" summary:"请求教师信息"`
	TeacherId int `json:"teacher_id" v:"required" dc:"教师id"`
}

type QFM struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"contents"`
	Views     int      `json:"views"`
	CreatedAt int64    `json:"created_at"`
	ImageURLs []string `json:"image_urls"`
}

type TeacherPinRes struct {
	QuestionList []QFM `json:"question_list"`
}
