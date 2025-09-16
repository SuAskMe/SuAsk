package notification

import (
	"context"
	"suask/internal/model"
	"suask/internal/service"
	"suask/utility/send_email"

	"github.com/gogf/gf/v2/frame/g"
)

func (s *sNotification) SendNoticeEmail(ctx context.Context, in *model.SendNoticeEmailInput) error {
	userSetting, err := service.Setting().GetSetting(ctx, model.GetSettingInput{Id: in.To})
	if err != nil {
		return err
	}
	if userSetting.NotifySwitch && userSetting.NotifyEmail != "" {
		er := send_email.SendNotice(userSetting.NotifyEmail, &send_email.Notice{
			User:    in.Notice.User,
			Type:    in.Notice.Type,
			Content: in.Notice.Content,
			URL:     in.Notice.URL,
		})
		if er != nil {
			g.Log("Email").Errorf(ctx, "send email to user %d error: %v", in.To, er)
		}
	}
	return nil
}
