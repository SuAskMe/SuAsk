// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// GuestUsers is the golang structure of table guest_users for DAO operations like Where/Data.
type GuestUsers struct {
	g.Meta    `orm:"table:guest_users, do:true"`
	Id        any         // 用户ID
	ExpiresAt *gtime.Time // 过期时间
	CreatedAt *gtime.Time // 创建时间
}
