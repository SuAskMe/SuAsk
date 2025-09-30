// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"suask/internal/model"
)

type (
	ITeacher interface {
		GetTeacherList(ctx context.Context, _ model.TeacherGetInput) (out *model.TeacherGetOutput, err error)
		GetTeacherAvatar(ctx context.Context, in *model.TeacherGetAvatarInput) (out *model.TeacherGetAvatarOutput, err error)
		TeacherExist(ctx context.Context, TeacherId int) (name string, err error)
		UpdatePerm(ctx context.Context, in model.TeacherUpdatePermInput) (out *model.TeacherUpdatePermOutput, err error)
	}
)

var (
	localTeacher ITeacher
)

func Teacher() ITeacher {
	if localTeacher == nil {
		panic("implement not found for interface ITeacher, forgot register?")
	}
	return localTeacher
}

func RegisterTeacher(i ITeacher) {
	localTeacher = i
}
