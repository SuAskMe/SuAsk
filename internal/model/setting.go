package model

type GetSettingInput struct {
	Id int `json:"id"         orm:"id"`
}

type GetSettingOutput struct {
	ThemeId int `json:"theme_id"    orm:"theme_id"`
	//Perm    interface{} `json:"perm" orm:"question_box_perm"`
}

type AddSettingInput struct {
	Id      int `json:"id"         orm:"id"`
	ThemeId int `json:"theme_id"    orm:"theme_id"`
	//Perm    interface{} `json:"perm" orm:"question_box_perm"`
}

type AddSettingOutput struct {
	Id int `json:"id"         orm:"id"`
}

type UpdateSettingInput struct {
	Id      int `json:"id"         orm:"id"`
	ThemeId int `json:"theme_id"    orm:"theme_id"`
	//Perm    interface{} `json:"perm" orm:"question_box_perm"`
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
