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
	IQuestionUtil interface {
		GetImages(ctx context.Context, input *model.GetImagesInput) (*model.GetImagesOutput, error)
		Favorite(ctx context.Context, input *model.FavoriteInput) (*model.FavoriteOutput, error)
	}
)

var (
	localQuestionUtil IQuestionUtil
)

func QuestionUtil() IQuestionUtil {
	if localQuestionUtil == nil {
		panic("implement not found for interface IQuestionUtil, forgot register?")
	}
	return localQuestionUtil
}

func RegisterQuestionUtil(i IQuestionUtil) {
	localQuestionUtil = i
}
