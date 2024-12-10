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
	IUser interface {
		GetUserInfoById(ctx context.Context, in model.GetUserInfoByIdInput) (out model.GetUserInfoByIdOutput, err error)
		GetUserInfo(ctx context.Context, in model.UserInfoInput) (out model.UserInfoOutput, err error)
		UpdateUserInfo(ctx context.Context, in model.UpdateUserInput) (out model.UpdateUserOutput, err error)
		UpdatePassword(ctx context.Context, in model.UpdatePasswordInput) (out model.UpdatePasswordOutput, err error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
