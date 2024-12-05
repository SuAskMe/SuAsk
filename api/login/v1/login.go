package v1

import "github.com/gogf/gf/v2/frame/g"

// 登录

type LoginReq struct {
	g.Meta   `path:"/user/login" method:"POST" tag:"Login" summary:"登录请求"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginRes struct {
	Token string `json:"token"`
}

// 登出

type LogoutReq struct {
	g.Meta `path:"/user/logout" method:"POST" tag:"Logout" summary:"登出请求"`
}

type LogoutRes struct{}

// 刷新 Token

type RefreshTokenReq struct {
	Token string `json:"token"`
}

type RefreshTokenRes struct {
	Token string `json:"token"`
}
