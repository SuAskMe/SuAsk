// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure for table questions.
type Questions struct {
	Id        int         `json:"id"        orm:"id"          description:""` //
	SrcUserId int         `json:"srcUserId" orm:"src_user_id" description:""` //
	DstUserId int         `json:"dstUserId" orm:"dst_user_id" description:""` //
	Title     string      `json:"title"     orm:"title"       description:""` //
	Contents  string      `json:"contents"  orm:"contents"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"  description:""` //
	Views     int         `json:"views"     orm:"views"       description:""` //
	ReplyCnt  int         `json:"replyCnt"  orm:"reply_cnt"   description:""` //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at"  description:""` //
}
