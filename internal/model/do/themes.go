// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Themes is the golang structure of table themes for DAO operations like Where/Data.
type Themes struct {
	g.Meta         `orm:"table:themes, do:true"`
	Id             interface{} // 主题ID
	BackgroundPath interface{} // 背景图片文件路径
}
