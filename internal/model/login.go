package model

type UserLoginInput struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
