package model

type Role string

type UserInfoBase struct {
	Id           int    `json:"id"           orm:"id"             description:"用户ID"`               // 用户ID
	Name         string `json:"name"         orm:"name"           description:"用户名"`                // 用户名
	Email        string `json:"email"        orm:"email"          description:"邮箱"`                 // 邮箱
	Nickname     string `json:"nickname"     orm:"nickname"       description:"昵称"`                 // 昵称
	Introduction string `json:"introduction" orm:"introduction"   description:"简介"`                 // 简介
	AvatarFileId int    `json:"avatarFileId" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"` // 头像文件ID，为空时为配置的默认头像
	ThemeId      int    `json:"themeId"      orm:"theme_id"       description:"主题ID，为空时为配置的默认主题"`   // 主题ID，为空时为配置的默认主题
}

const (
	Teacher Role = "teacher"
	Student Role = "student"
	Admin   Role = "admin"
)

// RegisterInput 注册输入
type RegisterInput struct {
	Name     string `json:"name" dc:"用户名"`
	Password string `json:"password" dc:"密码"`
	Role     Role   `json:"role" dc:"角色"`
	Token    string `json:"token" dc:"注册邮箱成功时传递的Token，用于在这里验证为同一个人，里面搭载 Email"`
}

// RegisterOutput 注册输出
type RegisterOutput struct {
	Name string `json:"name"`
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
	Id           int    `json:"id" v:"required" orm:"id" description:"用户ID"`
	Nickname     string `json:"nickname" v:"required" orm:"nickname" description:"昵称"`
	Introduction string `json:"introduction" v:"required" orm:"introduction" description:"简介"`
	AvatarFileId int    `json:"avatarFileId" v:"required" orm:"avatar_file_id" description:"头像文件ID，为空时为配置的默认头像"`
	ThemeId      int    `json:"themeId" v:"required" orm:"theme_id" description:"主题ID，为空时为配置的默认主题"`
}

type UpdateUserOutput struct {
}

type UpdatePasswordInput struct {
	Password string `json:"password" v:"required" dc:"新的密码"`
}

type UpdatePasswordOutput struct {
}
