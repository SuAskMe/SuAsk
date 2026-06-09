package model

type UserLoginInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginOutput struct {
	Role string `json:"role" orm:"role" dc:"用户角色"`
	Id   int    `json:"id"   orm:"id"   dc:"用户ID"`
}
