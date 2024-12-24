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
		GetTeacher(ctx context.Context, _ model.TeacherGetInput) (out model.TeacherGetOutput, err error)
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
