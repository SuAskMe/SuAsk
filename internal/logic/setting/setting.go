package setting

import (
	"context"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type sSetting struct {
}

func (s *sSetting) GetSetting(ctx context.Context, in model.GetSettingInput) (out model.GetSettingOutput, err error) {
	setting := model.GetSettingOutput{}
	// fmt.Println("id", in.Id)
	err = dao.Settings.Ctx(ctx).WherePri(in.Id).Scan(&setting)
	if err != nil {
		return model.GetSettingOutput{}, err
	}
	return setting, nil
}

func (s *sSetting) AddSetting(ctx context.Context, in model.AddSettingInput) (out model.AddSettingOutput, err error) {
	input := do.Settings{
		Id:           in.Id,
		ThemeId:      in.ThemeId,
		NotifySwitch: in.NotifySwitch,
		NotifyEmail:  in.NotifyEmail,
	}
	id, err := dao.Settings.Ctx(ctx).InsertAndGetId(input)
	if err != nil {
		return model.AddSettingOutput{}, err
	}
	return model.AddSettingOutput{
		Id: gconv.Int(id),
	}, nil
}

func (s *sSetting) UpdateSetting(ctx context.Context, in model.UpdateSettingInput) (out model.UpdateSettingOutput, err error) {
	input := do.Settings{
		Id:           in.Id,
		ThemeId:      in.ThemeId,
		NotifySwitch: in.NotifySwitch,
		NotifyEmail:  in.NotifyEmail,
	}
	_, err = dao.Settings.Ctx(ctx).WherePri(in.Id).Update(input)
	if err != nil {
		return model.UpdateSettingOutput{}, err
	}
	return model.UpdateSettingOutput{
		Id: in.Id,
	}, err
}

func init() {
	service.RegisterSetting(New())
}

func New() *sSetting {
	return &sSetting{}
}
