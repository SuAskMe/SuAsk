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
	ISetting interface {
		GetSetting(ctx context.Context, in model.GetSettingInput) (out model.GetSettingOutput, err error)
		AddSetting(ctx context.Context, in model.AddSettingInput) (out model.AddSettingOutput, err error)
		UpdateSetting(ctx context.Context, in model.UpdateSettingInput) (out model.UpdateSettingOutput, err error)
	}
)

var (
	localSetting ISetting
)

func Setting() ISetting {
	if localSetting == nil {
		panic("implement not found for interface ISetting, forgot register?")
	}
	return localSetting
}

func RegisterSetting(i ISetting) {
	localSetting = i
}
