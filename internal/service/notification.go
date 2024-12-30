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
	INotification interface {
		Add(ctx context.Context, in model.AddNotificationInput) (out model.AddNotificationOutput, err error)
		Get(ctx context.Context, in model.GetNotificationsInput) (out model.GetNotificationsOutput, err error)
		Update(ctx context.Context, in model.UpdateNotificationInput) (out model.UpdateNotificationOutput, err error)
		Delete(ctx context.Context, in model.DeleteNotificationInput) (out model.DeleteNotificationOutput, err error)
		NewNotificationCount(ctx context.Context, in model.NewNotificationCountInput) (out model.NewNotificationCountOutput, err error)
	}
)

var (
	localNotification INotification
)

func Notification() INotification {
	if localNotification == nil {
		panic("implement not found for interface INotification, forgot register?")
	}
	return localNotification
}

func RegisterNotification(i INotification) {
	localNotification = i
}
