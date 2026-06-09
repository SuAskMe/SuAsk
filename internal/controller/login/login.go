package login

import (
	"context"
	v1 "suask/api/login/v1"
	"suask/internal/consts"
	"suask/internal/middleware"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/session"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type cLogin struct{}

var Login cLogin

func (c *cLogin) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := g.RequestFromCtx(ctx)

	out, err := service.Login().Login(ctx, &model.UserLoginInput{Email: req.Email, Name: req.Name, Password: req.Password})
	if err != nil {
		return nil, err
	}

	// Create session and set cookie.
	sid, err := session.CreateSession(ctx, out.Id, out.Role)
	if err != nil {
		return nil, err
	}
	middleware.SetSessionCookie(r, sid)

	res = &v1.LoginRes{
		Role: out.Role,
		Id:   out.Id,
	}
	return
}

func (c *cLogin) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := gconv.Int(ctx.Value(consts.CtxId))

	// Clear the session from Redis.
	cookie := r.Cookie.Get(middleware.CookieName())
	if cookie != nil && cookie.String() != "" {
		_ = session.DeleteSession(ctx, cookie.String())
	}
	// Also clean up any user session index.
	_ = session.DeleteUserSessions(ctx, userId)

	middleware.ClearSessionCookie(r)
	return &v1.LogoutRes{UserId: userId}, nil
}

func (c *cLogin) HeartBeats(ctx context.Context, req *v1.HeartBeatsReq) (res *v1.HeartBeatsRes, err error) {
	r := g.RequestFromCtx(ctx)
	userId := gconv.Int(ctx.Value(consts.CtxId))

	// Refresh session TTL via the existing session middleware background refresh.
	// Also update guest_users expires_at if guest.
	cookie := r.Cookie.Get(middleware.CookieName())
	if cookie != nil && cookie.String() != "" {
		_, _ = session.RefreshSessionTTL(ctx, cookie.String())
	}

	_ = service.Login().HeartBeats(ctx)

	return &v1.HeartBeatsRes{UserId: userId}, nil
}
