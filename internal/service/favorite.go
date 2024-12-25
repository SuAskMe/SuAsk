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
		GetBase(ctx context.Context, in *model.GetFavoriteBaseInput) (out *model.GetFavoriteBaseOutput, err error)
		GetKeyWord(ctx context.Context, in *model.GetFavoriteKeywordsInput) (out *model.GetFavoriteKeywordsOutput, err error)
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
