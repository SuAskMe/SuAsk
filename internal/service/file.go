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
	IFile interface {
		UploadFile(ctx context.Context, in model.FileUploadInput) (out *model.FileUploadOutput, err error)
		Get(ctx context.Context, in model.FileGetInput) (out model.FileGetOutput, err error)
		GetList(ctx context.Context, in model.FileListGetInput) (out model.FileListGetOutput, err error)
	}
)

var (
	localFile IFile
)

func File() IFile {
	if localFile == nil {
		panic("implement not found for interface IFile, forgot register?")
	}
	return localFile
}

func RegisterFile(i IFile) {
	localFile = i
}
