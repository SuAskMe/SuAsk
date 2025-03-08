// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure for table questions.
type Questions struct {
	Id        int         `json:"id"        orm:"id"          description:"问题ID"`                    // 问题ID
	SrcUserId int         `json:"srcUserId" orm:"src_user_id" description:"发起提问的用户ID"`               // 发起提问的用户ID
	DstUserId int         `json:"dstUserId" orm:"dst_user_id" description:"被提问的用户ID，为空时问大家，不为空时问教师"` // 被提问的用户ID，为空时问大家，不为空时问教师
	Title     string      `json:"title"     orm:"title"       description:"问题标题"`                    // 问题标题
	Contents  string      `json:"contents"  orm:"contents"    description:"问题内容"`                    // 问题内容
	IsPrivate bool        `json:"isPrivate" orm:"is_private"  description:"是否私密提问，仅在问教师时可为是"`        // 是否私密提问，仅在问教师时可为是
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"  description:"创建时间"`                    // 创建时间
	Views     int         `json:"views"     orm:"views"       description:"浏览量"`                     // 浏览量
	ReplyCnt  int         `json:"replyCnt"  orm:"reply_cnt"   description:"回复数"`                     // 回复数
}
