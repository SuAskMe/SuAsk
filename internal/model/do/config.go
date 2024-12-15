// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Config is the golang structure of table config for DAO operations like Where/Data.
type Config struct {
	g.Meta            `orm:"table:config, do:true"`
	Id                interface{} // 配置ID，限制为0
	DefaultAvatarPath interface{} // 默认头像文件路径
	DefaultThemeId    interface{} // 默认主题ID
}
