package model

type Role string

// RegisterInput 注册输入
type RegisterInput struct {
	Name     string `json:"name" orm:"name" dc:"用户名"`
	UserSalt string `json:"userSalt" orm:"salt" dc:"加密盐"`
	Password string `json:"password" orm:"password" dc:"密码"`
	Role     string `json:"role" orm:"role" dc:"角色"`
	Email    string `json:"email" orm:"email" dc:"注册邮箱"`
	Token    string `json:"token" dc:"注册邮箱成功时传递的Token，用于在这里验证为同一个人，里面搭载 Email"`
}

// RegisterOutput 注册输出
type RegisterOutput struct {
	Id int `json:"id"`
}

// SendVerificationCodeInput 发送验证码输入
type SendVerificationCodeInput struct {
	Email string `json:"email" v:"required" dc:"要发送的邮箱地址"`
}

// SendVerificationCodeOutput 发送验证码输出
type SendVerificationCodeOutput struct {
	VerificationCode string `json:"verification_code"`
}

// VerifyVerificationCodeInput 验证验证码输入
type VerifyVerificationCodeInput struct {
	Email            string `json:"email" v:"required" dc:"要验证的邮箱"`
	VerificationCode string `json:"verification_code" v:"required" dc:"要验证的验证码"`
}

// VerifyVerificationCodeOutput 验证验证码输出
type VerifyVerificationCodeOutput struct {
	Token string `json:"token" dc:"验证成功的Token"`
}

type UpdateUserInput struct {
	Nickname     string `json:"nickname" orm:"nickname" dc:"昵称"`
	Introduction string `json:"introduction" orm:"introduction" dc:"简介"`
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id" dc:"头像文件ID，为空时为配置的默认头像"`
	ThemeId      int    `json:"themeId" orm:"theme_id" dc:"主题ID，为空时为配置的默认主题"`
}

type UpdateUserOutput struct {
	Id int `json:"id" v:"required" orm:"id" dc:"用户ID"`
}

type UpdatePasswordInput struct {
	Password string `json:"password" v:"required" dc:"新的密码"`
	Salt     string `json:"salt" v:"required" orm:"salt"`
}

type UpdatePasswordOutput struct {
	Id int `json:"id" v:"required" dc:"用户ID"`
}

type GetUserInfoByIdInput struct {
	Id int `json:"id" v:"required" dc:"用户ID"`
}

type GetUserInfoByIdOutput struct {
	UserInfoBase
}

type UserInfoInput struct {
	Id int `json:"id" v:"required" dc:"用户ID"`
}
type UserInfoOutput struct {
	UserInfoBase
	Email   string `json:"email"        orm:"email"          description:"邮箱"`
	ThemeId int    `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`
}

type UserInfoBase struct {
	Id           int    `json:"id"           orm:"id"             description:"用户ID"`
	Name         string `json:"name"         orm:"name"           description:"用户名"`
	Role         string `json:"role"         orm:"role"           description:"角色"`
	Nickname     string `json:"nickname"     orm:"nickname"       description:"昵称"`
	Introduction string `json:"introduction" orm:"introduction"   description:"简介"`
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"`
}
