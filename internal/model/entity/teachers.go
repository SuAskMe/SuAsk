// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Teachers is the golang structure for table teachers.
type Teachers struct {
	Id        int    `json:"id"        orm:"id"        description:""`
	Perm      string `json:"perm"      orm:"perm"      description:"提问箱权限"`
	Responses int    `json:"responses" orm:"responses" description:"回复数"`
	AvatarUrl string `json:"avatarUrl" orm:"avatar_url" description:"老师头像链接（deprecated）"`
}
