// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Settings is the golang structure of table settings for DAO operations like Where/Data.
type Settings struct {
	g.Meta          `orm:"table:settings, do:true"`
	Id              any // 设置id
	ThemeId         any // 主题ID，为空时为配置的默认主题
	QuestionBoxPerm any // 提问箱权限
	NotifyEmail     any // 通知邮箱
	NotifyMergeCnt  any // 合并提醒数量
	NotifyMaxDelay  any // 最大延迟时间(min)
	NotifySwitch    any // 通知开关
}
