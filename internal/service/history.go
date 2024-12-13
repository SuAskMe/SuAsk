// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import(
	"context"
	"suask/internal/model"
)

type (
	IHistory interface {
		LoadHistoryInfo(ctx context.Context, in *model.GetHistoryInput) (out *model.GetHistoryOutput, err error)
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
