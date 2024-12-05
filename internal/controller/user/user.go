package user

import (
	"context"
	v1 "suask/api/user/v1"
)

type cUser struct {
}

var User cUser

func (c *cUser) UpdateUserInfo(ctx context.Context, req v1.UpdateUserReq) (res v1.UpdateUserRes, err error) {
	return
}

func (c *cUser) UpdatePassWord(ctx context.Context, req v1.UpdatePasswordReq) (res v1.UpdatePasswordRes, err error) {
	return
}
