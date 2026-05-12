package middleware

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"suask/internal/consts"
	"suask/module/sjwt"
	"suask/utility/resp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const MsgErrAuthHeader = "Authorization : %s get token key fail"
const MsgErrAuthJwt = "Authorization : %s validate token fail"

// JwtRequired 中间件：必须登录，无有效 token 则返回 401。
func JwtRequired(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	claims, err := jwtAuth(r.Context(), authHeader)
	if err != nil {
		g.Log().Error(r.Context(), errors.Join(err, errors.New("must login")))
		resp.Do(r, 401, "请登录", nil)
		return
	}
	r.SetCtxVar(consts.CtxId, claims.UserID)
	g.Log().Debug(r.Context(), "URL", r.URL.Path, "UserId", claims.UserID)
	r.Middleware.Next()
}

// JwtOptional 中间件：可选登录，有 token 则解析用户 ID，无 token 则设为默认用户。
func JwtOptional(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	claims, err := jwtAuth(r.Context(), authHeader)
	if err != nil {
		r.SetCtxVar(consts.CtxId, consts.DefaultUserId)
	} else {
		r.SetCtxVar(consts.CtxId, claims.UserID)
	}
	g.Log().Debug(r.Context(), "URL", r.URL.Path, "Claims", claims)
	r.Middleware.Next()
}

// jwtAuth 解析并验证 JWT token。
func jwtAuth(ctx context.Context, authHeader string) (claims *sjwt.JwtClaims, err error) {
	if len(authHeader) == 0 {
		err = errors.New(sjwt.MsgLog(MsgErrAuthHeader, authHeader))
		return
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		err = errors.New(sjwt.MsgLog(MsgErrAuthHeader, authHeader))
		return
	} else if parts[1] == "" {
		err = errors.New(sjwt.MsgLog(MsgErrAuthHeader, authHeader))
		return
	}
	claims, err = sjwt.ParseToken(parts[1])
	if err != nil || claims == nil {
		err = errors.New(sjwt.MsgLog(MsgErrAuthJwt, authHeader))
		return
	}
	v, err := g.Redis().Get(ctx, consts.RedisJWTPrefix+strconv.Itoa(claims.UserID))
	if err != nil {
		return
	}
	if v.String() != parts[1] {
		err = errors.New(sjwt.MsgLog(MsgErrAuthJwt, authHeader))
		return
	}
	return
}
