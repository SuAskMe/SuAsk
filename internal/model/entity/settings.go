// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Settings is the golang structure for table settings.
type Settings struct {
	Id              int    `json:"id"              orm:"id"                description:"设置id"`             // 设置id
	ThemeId         int    `json:"themeId"         orm:"theme_id"          description:"主题ID，为空时为配置的默认主题"` // 主题ID，为空时为配置的默认主题
	QuestionBoxPerm string `json:"questionBoxPerm" orm:"question_box_perm" description:"提问箱权限"`            // 提问箱权限
	NotifyEmail     string `json:"notifyEmail"     orm:"notify_email"      description:"通知邮箱"`             // 通知邮箱
	NotifyMergeCnt  int    `json:"notifyMergeCnt"  orm:"notify_merge_cnt"  description:"合并提醒数量"`           // 合并提醒数量
	NotifyMaxDelay  int    `json:"notifyMaxDelay"  orm:"notify_max_delay"  description:"最大延迟时间(min)"`      // 最大延迟时间(min)
	NotifySwitch    bool   `json:"notifySwitch"    orm:"notify_switch"     description:"通知开关"`             // 通知开关
}
