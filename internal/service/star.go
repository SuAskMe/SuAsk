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
	IStar interface {
		GetStar(ctx context.Context, id int) (out model.StarQuestionOutPut, err error)
		DeleteStar(ctx context.Context, in model.DeleteStarInput) (out model.DeleteStarOutput, err error)
	}
)

var (
	localStar IStar
)

func Star() IStar {
	if localStar == nil {
		panic("implement not found for interface IStar, forgot register?")
	}
	return localStar
}

func RegisterStar(i IStar) {
	localStar = i
}
