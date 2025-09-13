package model

type GetSettingInput struct {
	Id int `json:"id"         orm:"id"`
}

type GetSettingOutput struct {
	ThemeId     int    `json:"theme_id"    orm:"theme_id"`
	NotifyEmail string `json:"notify_email"    orm:"notify_email"`
	//Perm    interface{} `json:"perm" orm:"question_box_perm"`
}

type AddSettingInput struct {
	Id          int    `json:"id"         orm:"id"`
	ThemeId     int    `json:"theme_id"    orm:"theme_id"`
	NotifyEmail string `json:"notify_email"    orm:"notify_email"`
	//Perm    interface{} `json:"perm" orm:"question_box_perm"`
}

type AddSettingOutput struct {
	Id int `json:"id"         orm:"id"`
}

type UpdateSettingInput struct {
	Id          int `json:"id"         orm:"id"`
	ThemeId     any `json:"theme_id"    orm:"theme_id"`
	NotifyEmail any `json:"notify_email"    orm:"notify_email"`
}

type UpdateSettingOutput struct {
	Id int `json:"id"         orm:"id"`
}

//type DeleteSettingInput struct {
//	Id int `json:"id"         orm:"id"`
//}
//
//type DeleteSettingOutput struct {
//	Id int `json:"id"         orm:"id"`
//}
