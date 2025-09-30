package model

type UserLoginInput struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type UserLoginOutput struct {
	Type  string `json:"type"         dc:"Token格式"`
	Token string `json:"token"        dc:"用户的Token"`
	Role  string `json:"role"         orm:"role"           dc:"用户角色"`
	Id    int    `json:"id"           orm:"id"             dc:"用户ID"`
}
