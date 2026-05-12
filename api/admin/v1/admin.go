package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ==================== 用户列表 ====================

type AdminUserItem struct {
	Id           int    `json:"id"            dc:"用户ID"`
	Name         string `json:"name"          dc:"用户名"`
	Nickname     string `json:"nickname"      dc:"昵称"`
	Email        string `json:"email"         dc:"邮箱"`
	Role         string `json:"role"          dc:"角色"`
	Introduction string `json:"introduction"  dc:"简介"`
	AvatarURL    string `json:"avatar"        dc:"头像URL"`
	CreatedAt    string `json:"created_at"    dc:"创建时间"`
}

type ListUsersReq struct {
	g.Meta  `path:"/admin/users" method:"GET" tags:"Admin" summary:"管理员-用户列表"`
	Page    int    `json:"page"    in:"query" v:"required|min:1" dc:"页码"`
	Role    string `json:"role"    in:"query" v:"in:admin,teacher,student" dc:"角色筛选"`
	Keyword string `json:"keyword" in:"query" dc:"搜索关键词"`
}

type ListUsersRes struct {
	List       []AdminUserItem `json:"list"`
	Total      int             `json:"total"`
	RemainPage int             `json:"remain_page"`
}

// ==================== 创建用户 ====================

type CreateUserReq struct {
	g.Meta       `path:"/admin/users" method:"POST" tags:"Admin" summary:"管理员-创建用户"`
	Name         string `json:"name"         v:"required|length:1,50" dc:"用户名"`
	Password     string `json:"password"     v:"required|length:6,72" dc:"密码"`
	Email        string `json:"email"        v:"required|email" dc:"邮箱"`
	Role         string `json:"role"         v:"required|in:admin,teacher,student" dc:"角色"`
	Nickname     string `json:"nickname"     v:"required|length:1,50" dc:"昵称"`
	Introduction string `json:"introduction" dc:"教师简介(仅teacher)"`
	Perm         string `json:"perm"         dc:"提问箱权限(仅teacher): public/protected/private"`
}

type CreateUserRes struct {
	Id int `json:"id" dc:"新用户ID"`
}

// ==================== 编辑用户 ====================

type UpdateUserReq struct {
	g.Meta       `path:"/admin/users/{id}" method:"PUT" tags:"Admin" summary:"管理员-编辑用户"`
	Id           int    `json:"id"           in:"path" v:"required" dc:"用户ID"`
	Nickname     string `json:"nickname"     dc:"昵称"`
	Email        string `json:"email"        dc:"邮箱"`
	Role         string `json:"role"         dc:"角色"`
	Introduction string `json:"introduction" dc:"教师简介"`
	Perm         string `json:"perm"         dc:"提问箱权限"`
}

type UpdateUserRes struct {
	Id int `json:"id" dc:"用户ID"`
}

// ==================== 修改头像 ====================

type UpdateAvatarReq struct {
	g.Meta `path:"/admin/users/{id}/avatar" method:"PUT" mime:"multipart/form-data" tags:"Admin" summary:"管理员-修改用户头像"`
	Id     int `json:"id" in:"path" v:"required" dc:"用户ID"`
}

type UpdateAvatarRes struct {
	Id        int    `json:"id" dc:"用户ID"`
	AvatarURL string `json:"avatar" dc:"新头像URL"`
}

// ==================== 重置密码 ====================

type ResetPasswordReq struct {
	g.Meta   `path:"/admin/users/{id}/password" method:"PUT" tags:"Admin" summary:"管理员-重置密码"`
	Id       int    `json:"id"       in:"path" v:"required" dc:"用户ID"`
	Password string `json:"password" v:"required|length:6,64" dc:"新密码"`
}

type ResetPasswordRes struct {
	Id int `json:"id" dc:"用户ID"`
}

// ==================== 删除用户 ====================

type DeleteUserReq struct {
	g.Meta `path:"/admin/users/{id}" method:"DELETE" tags:"Admin" summary:"管理员-删除用户"`
	Id     int `json:"id" in:"path" v:"required" dc:"用户ID"`
}

type DeleteUserRes struct {
	Id int `json:"id" dc:"用户ID"`
}
