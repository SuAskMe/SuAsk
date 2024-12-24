package notification

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
)

type sNotification struct{}

func (s *sNotification) Add(ctx context.Context, in model.AddNotificationInput) (out model.AddNotificationOutput, err error) {
	notification := do.Notifications{}
	err = gconv.Scan(in, &notification)
	if err != nil {
		return model.AddNotificationOutput{}, err
	}
	id, err := dao.Notifications.Ctx(ctx).InsertAndGetId(notification)
	if err != nil {
		return model.AddNotificationOutput{}, err
	}
	out = model.AddNotificationOutput{
		Id: int(id),
	}
	return out, nil
}

func (s *sNotification) Get(ctx context.Context, in model.GetNotificationsInput) (out model.GetNotificationsOutput, err error) {
	var notificationList []entity.Notifications
	var count int
	err = dao.Notifications.Ctx(ctx).Where(dao.Notifications.Columns().UserId, in.UserId).ScanAndCount(&notificationList, &count, false)
	if err != nil {
		return out, err
	}
	out = model.GetNotificationsOutput{
		Notifications: make([]model.Notification, count),
	}
	for i := range notificationList {
		out.Notifications[i] = model.Notification{
			Id:         notificationList[i].Id,
			QuestionId: notificationList[i].QuestionId,
			AnswerId:   notificationList[i].AnswerId,
			Type:       notificationList[i].Type,
			IsRead:     notificationList[i].IsRead,
			CreatedAt:  notificationList[i].CreatedAt,
		}
	}
	return out, nil
}

func (s *sNotification) Update(ctx context.Context, in model.UpdateNotificationInput) (out model.UpdateNotificationOutput, err error) {
	out = model.UpdateNotificationOutput{}
	_, err = dao.Notifications.Ctx(ctx).WherePri(in.Id).Update(do.Notifications{IsRead: true})
	if err != nil {
		return model.UpdateNotificationOutput{}, err
	}
	out.IsRead = true
	out.Id = in.Id
	return out, nil
}

func (s *sNotification) Delete(ctx context.Context, in model.DeleteNotificationInput) (out model.DeleteNotificationOutput, err error) {
	_, err = dao.Notifications.Ctx(ctx).Delete(in.Id)
	if err != nil {
		return model.DeleteNotificationOutput{}, err
	}
	out = model.DeleteNotificationOutput{}
	return out, nil
}

func init() {
	service.RegisterNotification(New())
}

func New() *sNotification {
	return &sNotification{}
}
