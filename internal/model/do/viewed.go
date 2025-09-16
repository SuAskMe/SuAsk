// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Viewed is the golang structure of table viewed for DAO operations like Where/Data.
type Viewed struct {
	g.Meta     `orm:"table:viewed, do:true"`
	UserId     any // 用户id
	QuestionId any // 问题id
}
