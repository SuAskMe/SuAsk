package middleware

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func CORS(r *ghttp.Request) {
	origin := r.Header.Get("Origin")
	if origin == "" {
		origin = g.Cfg().MustGet(r.Context(), "server.corsOrigin", "*").String()
	}

	r.Response.CORS(ghttp.CORSOptions{
		AllowOrigin:      origin,
		AllowCredentials: "true",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type,Authorization,X-Admin-Mode,X-Device-Id",
	})

	if r.Method == "OPTIONS" {
		r.ExitAll()
		return
	}

	r.Middleware.Next()
}
