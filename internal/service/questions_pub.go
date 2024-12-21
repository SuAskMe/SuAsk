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
	IPublicQuestion interface {
		GetBase(ctx context.Context, input *model.GetBaseInput) (*model.GetBaseOutput, error)
		GetAnswers(ctx context.Context, input *model.GetAnswersInput) (*model.GetAnswersOutput, error)
		GetKeyword(ctx context.Context, input *model.GetKeywordsInput) (*model.GetKeywordsOutput, error)
		AddQuestion(ctx context.Context, in *model.AddQuestionInput) (out *model.AddQuestionOutput, err error)
	}
)

var (
	localPublicQuestion IPublicQuestion
)

func PublicQuestion() IPublicQuestion {
	if localPublicQuestion == nil {
		panic("implement not found for interface IPublicQuestion, forgot register?")
	}
	return localPublicQuestion
}

func RegisterPublicQuestion(i IPublicQuestion) {
	localPublicQuestion = i
}
