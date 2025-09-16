package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
	"time"
)

func Middleware(token *gtoken.GfToken, ctx context.Context, group *ghttp.RouterGroup) error {
	m := token
	if !m.InitConfig() {
		return errors.New("InitConfig fail")
	}

	// 设置为Group模式
	m.MiddlewareType = gtoken.MiddlewareTypeGroup
	g.Log().Info(ctx, "[GToken][params:"+m.String()+"]start... ")

	// 缓存模式
	if m.CacheMode > gtoken.CacheModeFile {
		g.Log().Error(ctx, "[GToken]CacheMode set error")
		return errors.New("CacheMode set error")
	}
	// 登录
	if m.LoginPath == "" || m.LoginBeforeFunc == nil {
		g.Log().Error(ctx, "[GToken]LoginPath or LoginBeforeFunc not set")
		return errors.New("LoginPath or LoginBeforeFunc not set")
	}
	// 登出
	if m.LogoutPath == "" {
		g.Log().Error(ctx, "[GToken]LogoutPath not set")
		return errors.New("LogoutPath not set")
	}

	group.Middleware(authMiddleware(m))

	registerFunc(ctx, group, m.LoginPath, m.Login)
	registerFunc(ctx, group, m.LogoutPath, m.Logout)

	return nil
}

func registerFunc(ctx context.Context, group *ghttp.RouterGroup, pattern string, object interface{}) {
	if gstr.Contains(pattern, ":") || gstr.Contains(pattern, "@") {
		group.Map(map[string]interface{}{
			pattern: object,
		})
	} else {
		group.ALL(pattern, object)
	}
}

// AuthMiddleware 认证拦截
func authMiddleware(m *gtoken.GfToken) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		urlPath := r.URL.Path
		if !m.AuthPath(r.Context(), urlPath) {
			// 如果不需要认证，继续
			r.Middleware.Next()
			return
		}
		// 不需要认证，直接下一步
		if !m.AuthBeforeFunc(r) {
			r.Middleware.Next()
			return
		}

		// 获取请求token
		tokenResp := getRequestToken(m, r)
		if tokenResp.Success() {
			// 验证token
			tokenResp = validToken(m, r.Context(), tokenResp.DataString())
		}
		m.AuthAfterFunc(r, tokenResp)
	}
}

// getRequestToken 返回请求Token
func getRequestToken(m *gtoken.GfToken, r *ghttp.Request) gtoken.Resp {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			g.Log().Warning(r.Context(), msgLog(gtoken.MsgErrAuthHeader, authHeader))
			return gtoken.Unauthorized(fmt.Sprintf(gtoken.MsgErrAuthHeader, authHeader), "")
		} else if parts[1] == "" {
			g.Log().Warning(r.Context(), msgLog(gtoken.MsgErrAuthHeader, authHeader))
			return gtoken.Unauthorized(fmt.Sprintf(gtoken.MsgErrAuthHeader, authHeader), "")
		}

		return gtoken.Succ(parts[1])
	}

	authHeader = r.Get(gtoken.KeyToken).String()
	if authHeader == "" {
		return gtoken.Unauthorized(gtoken.MsgErrTokenEmpty, "")
	}
	return gtoken.Succ(authHeader)

}

// validToken 验证Token
func validToken(m *gtoken.GfToken, ctx context.Context, token string) gtoken.Resp {
	if token == "" {
		return gtoken.Unauthorized(gtoken.MsgErrTokenEmpty, "")
	}

	decryptToken := m.DecryptToken(ctx, token)
	if !decryptToken.Success() {
		return decryptToken
	}

	userKey := decryptToken.GetString(gtoken.KeyUserKey)
	uuid := decryptToken.GetString(gtoken.KeyUuid)

	userCacheResp := getToken(m, ctx, userKey)
	if !userCacheResp.Success() {
		return userCacheResp
	}

	if uuid != userCacheResp.GetString(gtoken.KeyUuid) {
		g.Log().Debug(ctx, msgLog(gtoken.MsgErrAuthUuid)+", decryptToken:"+decryptToken.Json()+" cacheValue:"+gconv.String(userCacheResp.Data))
		return gtoken.Unauthorized(gtoken.MsgErrAuthUuid, "")
	}

	return userCacheResp
}

func msgLog(msg string, params ...interface{}) string {
	if len(params) == 0 {
		return gtoken.DefaultLogPrefix + msg
	}
	return gtoken.DefaultLogPrefix + fmt.Sprintf(msg, params...)
}

// getToken 通过userKey获取Token
func getToken(m *gtoken.GfToken, ctx context.Context, userKey string) gtoken.Resp {
	cacheKey := m.CacheKey + userKey

	userCacheResp := getCache(m, ctx, cacheKey)
	if !userCacheResp.Success() {
		return userCacheResp
	}
	userCache := gconv.Map(userCacheResp.Data)

	nowTime := gtime.Now().TimestampMilli()
	refreshTime := userCache[gtoken.KeyRefreshTime]

	// 需要进行缓存超时时间刷新
	if gconv.Int64(refreshTime) == 0 || nowTime > gconv.Int64(refreshTime) {
		userCache[gtoken.KeyCreateTime] = gtime.Now().TimestampMilli()
		userCache[gtoken.KeyRefreshTime] = gtime.Now().TimestampMilli() + gconv.Int64(m.MaxRefresh)
		return setCache(m, ctx, cacheKey, userCache)
	}

	return gtoken.Succ(userCache)
}

// setCache 设置缓存
func setCache(m *gtoken.GfToken, ctx context.Context, cacheKey string, userCache g.Map) gtoken.Resp {
	switch m.CacheMode {
	case gtoken.CacheModeCache, gtoken.CacheModeFile:
		gcache.Set(ctx, cacheKey, userCache, gconv.Duration(m.Timeout)*time.Millisecond)
		if m.CacheMode == gtoken.CacheModeFile {
			writeFileCache(ctx)
		}
	case gtoken.CacheModeRedis:
		cacheValueJson, err1 := gjson.Encode(userCache)
		if err1 != nil {
			g.Log().Error(ctx, "[GToken]cache json encode error", err1)
			return gtoken.Error("cache json encode error")
		}
		_, err := g.Redis().Do(ctx, "SETEX", cacheKey, m.Timeout/1000, cacheValueJson)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache set error", err)
			return gtoken.Error("cache set error")
		}
	default:
		return gtoken.Error("cache model error")
	}

	return gtoken.Succ(userCache)
}

// getCache 获取缓存
func getCache(m *gtoken.GfToken, ctx context.Context, cacheKey string) gtoken.Resp {
	var userCache g.Map
	switch m.CacheMode {
	case gtoken.CacheModeCache, gtoken.CacheModeFile:
		userCacheValue, err := gcache.Get(ctx, cacheKey)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache get error", err)
			return gtoken.Error("cache get error")
		}
		if userCacheValue.IsNil() {
			return gtoken.Unauthorized("login timeout or not login", "")
		}
		userCache = gconv.Map(userCacheValue)
	case gtoken.CacheModeRedis:
		userCacheJson, err := g.Redis().Do(ctx, "GET", cacheKey)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache get error", err)
			return gtoken.Error("cache get error")
		}
		if userCacheJson.IsNil() {
			return gtoken.Unauthorized("login timeout or not login", "")
		}

		err = gjson.DecodeTo(userCacheJson, &userCache)
		if err != nil {
			g.Log().Error(ctx, "[GToken]cache get json error", err)
			return gtoken.Error("cache get json error")
		}
	default:
		return gtoken.Error("cache model error")
	}

	return gtoken.Succ(userCache)
}

func writeFileCache(ctx context.Context) {
	file := gfile.Temp(gtoken.CacheModeFileDat)
	data, e := gcache.Data(ctx)
	if e != nil {
		g.Log().Error(ctx, "[GToken]cache writeFileCache error", e)
	}
	gfile.PutContents(file, gjson.New(data).MustToJsonString())
}
