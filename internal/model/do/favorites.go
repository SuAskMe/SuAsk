// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Favorites is the golang structure of table favorites for DAO operations like Where/Data.
type Favorites struct {
	g.Meta     `orm:"table:favorites, do:true"`
	Id         any         // 收藏（置顶）ID
	UserId     any         // 用户ID
	QuestionId any         // 问题ID
	CreatedAt  *gtime.Time // 创建时间
	Package    any         // 收藏夹
}
