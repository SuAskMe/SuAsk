package login

import (
	"context"
	"strconv"
	v1 "suask/api/login/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type cLogin struct{}

var Login cLogin

func (c *cLogin) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	out, err := service.Login().Login(ctx, &model.UserLoginInput{Email: req.Email, Name: req.Name, Password: req.Password})
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(out, &res)
	if err != nil {
		return nil, err
	}
	return
}

func (c *cLogin) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	g.Redis().Del(ctx, consts.RedisJWTPrefix+strconv.Itoa(userId))
	return
}

func (c *cLogin) HeartBeats(ctx context.Context, req *v1.HeartBeatsReq) (res *v1.HeartBeatsRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	return &v1.HeartBeatsRes{UserId: userId}, nil
}

// func (c *cLogin) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
// 	return
// }
