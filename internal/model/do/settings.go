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
	Id              interface{} // 设置id
	ThemeId         interface{} // 主题ID，为空时为配置的默认主题
	QuestionBoxPerm interface{} // 提问箱权限
}
