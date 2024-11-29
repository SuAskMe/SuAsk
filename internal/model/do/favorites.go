// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Favorites is the golang structure of table favorites for DAO operations like Where/Data.
type Favorites struct {
	g.Meta     `orm:"table:favorites, do:true"`
	Id         interface{} // 收藏（置顶）ID
	UserId     interface{} // 用户ID
	QuestionId interface{} // 问题ID
	Ordinal    interface{} // 收藏序号
}
