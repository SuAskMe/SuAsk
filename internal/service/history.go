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
	IHistoryOperation interface {
		// 查找历史提问模块需要的信息
		LoadHistoryInfo(ctx context.Context, in *model.GetHistoryInput) (out *model.GetHistoryOutput, err error)
	}
)

var (
	localHistoryOperation IHistoryOperation
)

func HistoryOperation() IHistoryOperation {
	if localHistoryOperation == nil {
		panic("implement not found for interface IHistoryOperation, forgot register?")
	}
	return localHistoryOperation
}

func RegisterHistoryOperation(i IHistoryOperation) {
	localHistoryOperation = i
}
