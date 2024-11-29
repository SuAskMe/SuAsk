// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Config is the golang structure for table config.
type Config struct {
	Id                  bool `json:"id"                  orm:"id"                     description:"配置ID，限制为0"` // 配置ID，限制为0
	DefaultAvatarFileId int  `json:"defaultAvatarFileId" orm:"default_avatar_file_id" description:"默认头像文件ID"`  // 默认头像文件ID
	DefaultThemeId      int  `json:"defaultThemeId"      orm:"default_theme_id"       description:"默认主题ID"`    // 默认主题ID
}
