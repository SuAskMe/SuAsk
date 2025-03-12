package setting

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
)

type sSetting struct {
}

func (s *sSetting) GetSetting(ctx context.Context, in model.GetSettingInput) (out model.GetSettingOutput, err error) {
	setting := model.GetSettingOutput{}
	fmt.Println("id", in.Id)
	err = dao.Settings.Ctx(ctx).WherePri(in.Id).Scan(&setting)
	if err != nil {
		return model.GetSettingOutput{}, err
	}
	return setting, nil
}

func (s *sSetting) AddSetting(ctx context.Context, in model.AddSettingInput) (out model.AddSettingOutput, err error) {
	input := do.Settings{
		Id:      in.Id,
		ThemeId: in.ThemeId,
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
		Id:      in.Id,
		ThemeId: in.ThemeId,
		//QuestionBoxPerm: in.Perm,
	}
	_, err = dao.Settings.Ctx(ctx).WherePri(in.Id).Update(input)
	if err != nil {
		return model.UpdateSettingOutput{}, err
	}
	return model.UpdateSettingOutput{
		Id: gconv.Int(in.Id),
	}, err
}

func init() {
	service.RegisterSetting(New())
}

func New() *sSetting {
	return &sSetting{}
}
