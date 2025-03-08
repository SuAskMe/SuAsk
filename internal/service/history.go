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
	IHistory interface {
		GetBase(ctx context.Context, in *model.GetHistoryBaseInput) (out *model.GetHistoryBaseOutput, err error)
		GetKeyWord(ctx context.Context, in *model.GetHistoryKeywordsInput) (out *model.GetHistoryKeywordsOutput, err error)
	}
)

var (
	localHistory IHistory
)

func History() IHistory {
	if localHistory == nil {
		panic("implement not found for interface IHistory, forgot register?")
	}
	return localHistory
}

func RegisterHistory(i IHistory) {
	localHistory = i
}
