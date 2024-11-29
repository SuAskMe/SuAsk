// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Config is the golang structure for table config.
type Config struct {
	Id                bool   `json:"id"                orm:"id"                  description:"配置ID，限制为0"` // 配置ID，限制为0
	DefaultAvatarPath string `json:"defaultAvatarPath" orm:"default_avatar_path" description:"默认头像文件路径"`  // 默认头像文件路径
	DefaultThemeId    int    `json:"defaultThemeId"    orm:"default_theme_id"    description:"默认主题ID"`    // 默认主题ID
}
