// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Teachers is the golang structure for table teachers.
type Teachers struct {
	Id           int    `json:"id"           orm:"id"           description:""`       //
	Responses    int    `json:"responses"    orm:"responses"    description:"回复数"`    // 回复数
	Name         string `json:"name"         orm:"name"         description:"老师名字"`   // 老师名字
	AvatarUrl    string `json:"avatarUrl"    orm:"avatar_url"   description:"老师头像链接"` // 老师头像链接
	Introduction string `json:"introduction" orm:"introduction" description:"老师简介"`   // 老师简介
	Email        string `json:"email"        orm:"email"        description:"老师邮箱"`   // 老师邮箱
	Perm         string `json:"perm"         orm:"perm"         description:"提问箱权限"`  // 提问箱权限
}
