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
	IFavorite interface {
		GetFavorite(ctx context.Context, id int) (out model.FavoriteQuestionOutPut, err error)
		GetPageFavorite(ctx context.Context, in model.PageFavoriteQuestionInPut) (out model.PageFavoriteQuestionOutPut, err error)
		DeleteFavorite(ctx context.Context, in model.DeleteFavoriteInput) (out model.DeleteFavoriteOutput, err error)
	}
)

var (
	localFavorite IFavorite
)

func Favorite() IFavorite {
	if localFavorite == nil {
		panic("implement not found for interface IFavorite, forgot register?")
	}
	return localFavorite
}

func RegisterFavorite(i IFavorite) {
	localFavorite = i
}
