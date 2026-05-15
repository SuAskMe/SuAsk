package middleware

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model/entity"
	"suask/utility/resp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// AdminRequired 中间件：验证当前用户为管理员角色，非管理员返回 403。
// 必须在 JwtRequired 之后使用，依赖其设置的 CtxId。
func AdminRequired(r *ghttp.Request) {
	userId := gconv.Int(r.Context().Value(consts.CtxId))

	var user entity.Users
	err := dao.Users.Ctx(r.Context()).
		Where(dao.Users.Columns().Id, userId).
		Fields(dao.Users.Columns().Role).
		Scan(&user)
	if err != nil {
		g.Log().Error(r.Context(), "AdminRequired: query user failed", "userId", userId, "err", err)
		resp.Do(r, 403, "需要管理员权限", nil)
		return
	}

	if user.Role != consts.ADMIN {
		resp.Do(r, 403, "需要管理员权限", nil)
		return
	}

	r.Middleware.Next()
}

// IsAdminMode 检查请求是否处于管理员模式。
// 当 X-Admin-Mode header 为 "true" 且当前用户角色为 admin 时返回 true。
// 非 admin 用户即使发送了 X-Admin-Mode header 也会被忽略（返回 false）。
func IsAdminMode(ctx context.Context) bool {
	r := g.RequestFromCtx(ctx)
	if r == nil {
		return false
	}

	adminMode := r.Header.Get("X-Admin-Mode")
	if adminMode != "true" {
		return false
	}

	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == 0 {
		return false
	}

	var user entity.Users
	err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userId).
		Fields(dao.Users.Columns().Role).
		Scan(&user)
	if err != nil {
		g.Log().Error(ctx, "IsAdminMode: query user failed", "userId", userId, "err", err)
		return false
	}

	return user.Role == consts.ADMIN
}
