// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// GuestUsers is the golang structure for table guest_users.
type GuestUsers struct {
	Id        int         `json:"id"        orm:"id"         description:"用户ID"`
	ExpiresAt *gtime.Time `json:"expiresAt" orm:"expires_at" description:"过期时间"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`
}
