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
	ITeacherQuestion interface {
		GetBase(ctx context.Context, input *model.GetBaseOfTeacherInput) (*model.GetBaseOfTeacherOutput, error)
		GetKeyword(ctx context.Context, input *model.GetKeywordsOfTeacherInput) (*model.GetKeywordsOutput, error)
	}
)

var (
	localTeacherQuestion ITeacherQuestion
)

func TeacherQuestion() ITeacherQuestion {
	if localTeacherQuestion == nil {
		panic("implement not found for interface ITeacherQuestion, forgot register?")
	}
	return localTeacherQuestion
}

func RegisterTeacherQuestion(i ITeacherQuestion) {
	localTeacherQuestion = i
}
