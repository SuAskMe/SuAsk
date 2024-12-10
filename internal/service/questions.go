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
		Get(ctx context.Context, input *model.GetPublicQuestionInput) (*model.GetPublicQuestionOutput, error)
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
