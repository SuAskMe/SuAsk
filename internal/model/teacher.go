package model

import v1 "suask/api/teacher/v1"

type TeacherGetInput struct {
}

type TeacherGetOutput struct {
	TeacherList []v1.TeacherBase
}

type TeacherGetAvatarInput struct {
	TeacherId int `json:"teacher_id" orm:"id"`
}

type TeacherGetAvatarOutput struct {
	AvatarUrl string `json:"avatar_url" orm:"avatar_url"`
}

type TeacherUpdatePermInput struct {
	TeacherId int    `json:"teacher_id" orm:"id"`
	Perm      string `json:"perm" v:"required|enums" orm:"perm" dc:"要更新的提问箱权限"`
}

type TeacherUpdatePermOutput struct {
	TeacherId int `json:"teacher_id" orm:"id" dc:"老师的id"`
}
