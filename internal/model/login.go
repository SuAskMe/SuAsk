package model

type UserLoginInput struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}
