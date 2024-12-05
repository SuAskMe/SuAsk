package login

import (
	"context"
	v1 "suask/api/login/v1"
)

type cLogin struct{}

var Login cLogin

func (c *cLogin) Login(ctx context.Context, req v1.LoginReq) (res v1.LoginRes, err error) {
	return
}

func (c *cLogin) Logout(ctx context.Context, req v1.LogoutReq) (res v1.LogoutRes, err error) {
	return
}

func (c *cLogin) RefreshToken(ctx context.Context, req v1.RefreshTokenReq) (res v1.RefreshTokenRes, err error) {
	return
}
