package model

type Role string

type UpdateUserInput struct {
	Nickname     interface{} `json:"nickname" orm:"nickname" dc:"昵称"`
	Introduction interface{} `json:"introduction" orm:"introduction" dc:"简介"`
	AvatarFileId interface{} `json:"avatarId" orm:"avatar_file_id" dc:"头像文件ID，为空时为配置的默认头像"`
	ThemeId      interface{} `json:"themeId" orm:"theme_id" dc:"主题ID，为空时为配置的默认主题"`
}

type UpdateUserOutput struct {
	Id int `json:"id" v:"required" orm:"id" dc:"用户ID"`
}

type UpdatePasswordInput struct {
	UserId   int    `json:"userId" dc:"用户ID"`
	Password string `json:"password" v:"required" dc:"新的密码"`
}

type UpdatePasswordOutput struct {
	Id int `json:"id" v:"required" dc:"用户ID"`
}

type UserInfoInput struct {
	Id int `json:"id" v:"required" dc:"用户ID"`
}
type UserInfoOutput struct {
	Id           int    `json:"id"           orm:"id"             description:"用户ID"`
	Name         string `json:"name"         orm:"name"           description:"用户名"`
	Role         string `json:"role"         orm:"role"           description:"角色"`
	Nickname     string `json:"nickname"     orm:"nickname"       description:"昵称"`
	Introduction string `json:"introduction" orm:"introduction"   description:"简介"`
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"`
	Email        string `json:"email"        orm:"email"          description:"邮箱"`
	ThemeId      int    `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`
}
