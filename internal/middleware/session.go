package middleware

import (
	"errors"
	"net/http"
	"time"

	"suask/internal/consts"
	"suask/module/session"
	"suask/utility/resp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// SessionRequired 中间件：必须登录，从 suask_sid cookie 验证会话。
func SessionRequired(r *ghttp.Request) {
	cookie := r.Cookie.Get(session.CookieName())
	if cookie == nil || cookie.Val() == "" {
		g.Log().Error(r.Context(), "SessionRequired: no session cookie")
		resp.Do(r, 401, "请登录", nil)
		return
	}

	sid := cookie.String()
	p, err := session.GetSession(r.Context(), sid)
	if err != nil {
		g.Log().Error(r.Context(), "SessionRequired: get session failed", "err", err, "sid", sid)
		resp.Do(r, 401, "请登录", nil)
		return
	}
	if p == nil {
		g.Log().Error(r.Context(), errors.New("SessionRequired: session not found"), "sid", sid)
		resp.Do(r, 401, "请登录", nil)
		return
	}

	r.SetCtxVar(consts.CtxId, p.UserID)
	r.SetCtxVar(consts.CtxRole, p.Role)
	g.Log().Debug(r.Context(), "URL", r.URL.Path, "UserId", p.UserID, "Role", p.Role)

	// Lightweight TTL refresh in background (do not block handler).
	go func(ctx *ghttp.Request, ssid string) {
		_, _ = session.RefreshSessionTTL(ctx.Context(), ssid)
	}(r, sid)

	r.Middleware.Next()
}

// SetSessionCookie writes an HttpOnly session cookie to the response.
func SetSessionCookie(r *ghttp.Request, sid string) {
	name, maxAge := session.CookieConfig()
	r.Cookie.SetCookie(name, sid, "", "/", time.Duration(maxAge)*time.Second, ghttp.CookieOptions{
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie expires the session cookie immediately.
func ClearSessionCookie(r *ghttp.Request) {
	name, _ := session.CookieConfig()
	r.Cookie.SetCookie(name, "", "", "/", 0, ghttp.CookieOptions{
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})
}

// CookieName returns the session cookie name.
func CookieName() string {
	return session.CookieName()
}
