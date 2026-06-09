package middleware

import (
	"context"
	"strings"
	"suask/internal/consts"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func CleanErrorResponse(r *ghttp.Request) {
	r.Middleware.Next()

	err := r.GetError()
	if err == nil {
		return
	}

	message := err.Error()
	if shouldSanitizeError(message) {
		g.Log().Error(r.Context(), "CleanErrorResponse: sanitize internal error", "err", err)
		r.SetError(gerror.New(consts.ErrInternal))
	}
}

func shouldSanitizeError(message string) bool {
	var internalMarkers = []string{
		"SQL",
		"sql:",
		"SELECT ",
		"INSERT ",
		"UPDATE ",
		"DELETE ",
		"Scan failed",
		"Count failed",
	}

	for _, marker := range internalMarkers {
		if strings.Contains(message, marker) {
			return true
		}
	}
	return false
}

func SanitizeError(ctx context.Context, err error, fallback string, logArgs ...any) error {
	if err == nil {
		return nil
	}
	args := append([]any{fallback}, logArgs...)
	args = append(args, "err", err)
	g.Log().Error(ctx, args...)
	return gerror.New(fallback)
}
