package notification

import (
	"context"
	v1 "suask/api/notification/v1"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cNotification struct{}

var Notification = cNotification{}

func (c *cNotification) Get(ctx context.Context, req *v1.NotificationGetReq) (res *v1.NotificationGetRes, err error) {
	in := model.GetNotificationsInput{UserId: req.UserId}
	out, err := service.Notification().Get(ctx, in)
	if err != nil {
		return nil, err
	}
	res = &v1.NotificationGetRes{}
	err = gconv.Scan(out, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *cNotification) Update(ctx context.Context, req *v1.NotificationUpdateReq) (res *v1.NotificationUpdateRes, err error) {
	in := model.UpdateNotificationInput{Id: req.Id}
	out, err := service.Notification().Update(ctx, in)
	if err != nil {
		return nil, err
	}
	res = &v1.NotificationUpdateRes{}
	err = gconv.Scan(out, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *cNotification) Delete(ctx context.Context, req *v1.NotificationDeleteReq) (res *v1.NotificationDeleteRes, err error) {
	in := model.DeleteNotificationInput{Id: req.Id}
	_, err = service.Notification().Delete(ctx, in)
	if err != nil {
		return nil, err
	}
	return &v1.NotificationDeleteRes{}, nil
}

func (c *cNotification) GetCount(ctx context.Context, req *v1.NotificationGetCountReq) (res *v1.NotificationGetCountRes, err error) {
	in := model.NewNotificationCountInput{UserId: req.UserId}
	out, err := service.Notification().NewNotificationCount(ctx, in)
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(out, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
