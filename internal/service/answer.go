// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"suask/internal/model/entity"
)

type (
	IAnswer interface {
		GetAnswerIDs(ctx context.Context, answerId int) (out entity.Answers, err error)
	}
)

var (
	localAnswer IAnswer
)

func Answer() IAnswer {
	if localAnswer == nil {
		panic("implement not found for interface IAnswer, forgot register?")
	}
	return localAnswer
}

func RegisterAnswer(i IAnswer) {
	localAnswer = i
}
