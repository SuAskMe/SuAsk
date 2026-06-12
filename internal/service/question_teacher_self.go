package service

import (
	"context"
	"suask/internal/model"
)

type (
	ITeacherQuestionSelf interface {
		GetQFMAll(ctx context.Context, input *model.GetQFMInput) (*model.GetQFMOutput, error)
		GetQFMPinned(ctx context.Context, input *model.GetQFMInput) (*model.GetQFMOutput, error)
		GetKeyword(ctx context.Context, input *model.GetQFMKeywordsInput) (*model.GetKeywordsOutput, error)
		PinQFM(ctx context.Context, input *model.PinQFMInput) (*model.PinQFMOutput, error)
	}
)

var localTeacherQuestionSelf ITeacherQuestionSelf

func TeacherQuestionSelf() ITeacherQuestionSelf {
	if localTeacherQuestionSelf == nil {
		panic("implement not found for interface ITeacherQuestionSelf, forgot register?")
	}
	return localTeacherQuestionSelf
}

func RegisterTeacherQuestionSelf(i ITeacherQuestionSelf) {
	localTeacherQuestionSelf = i
}
