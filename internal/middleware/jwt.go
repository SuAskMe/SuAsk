package middleware

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"suask/internal/consts"
	"suask/module/sjwt"
	"suask/utility/resp"
	triemux "suask/utility/trie_mux"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

const MsgErrAuthHeader = "Authorization : %s get token key fail"
const MsgErrAuthJwt = "Authorization : %s validate token fail"

type JWTMiddleware struct {
	trieMux        *triemux.TrieMux
	excpetUrlMatch []string
	fullUrlMatch   []string
}

func (j *JWTMiddleware) BuildMustLoginTrie(prefixUrlMatch, fullUrlMatch, excpetUrlMatch []string) {
	j.trieMux = triemux.NewTrieMux()
	for _, v := range prefixUrlMatch {
		if err := j.trieMux.Insert(v); err != nil {
			panic(err)
		}
	}
	j.excpetUrlMatch = excpetUrlMatch
	j.fullUrlMatch = fullUrlMatch
}

func (j *JWTMiddleware) JwtAuth(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	claims, err := j.auth(r.Context(), authHeader)
	if err != nil {
		if j.isMustLoginPath(r.URL.String()) {
			g.Log().Error(r.Context(), errors.Join(err, errors.New("is must login path")))
			resp.Do(r, 401, "请登录", nil)
			return
		}
		r.SetCtxVar(consts.CtxId, consts.DefaultUserId)
	} else {
		r.SetCtxVar(consts.CtxId, claims.UserID)
	}
	r.Middleware.Next()
}

func (j *JWTMiddleware) auth(ctx context.Context, authHeader string) (claims *sjwt.JwtClaims, err error) {
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

func (j *JWTMiddleware) isMustLoginPath(path string) bool {
	if slices.Contains(j.excpetUrlMatch, path) {
		return false
	}
	if j.trieMux.HasPrefix(path) {
		return true
	}
	return slices.Contains(j.fullUrlMatch, path)
}
