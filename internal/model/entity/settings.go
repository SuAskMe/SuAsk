// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Settings is the golang structure for table settings.
type Settings struct {
	Id              int    `json:"id"              orm:"id"                description:"设置id"`             // 设置id
	ThemeId         int    `json:"themeId"         orm:"theme_id"          description:"主题ID，为空时为配置的默认主题"` // 主题ID，为空时为配置的默认主题
	QuestionBoxPerm string `json:"questionBoxPerm" orm:"question_box_perm" description:"提问箱权限"`            // 提问箱权限
}
