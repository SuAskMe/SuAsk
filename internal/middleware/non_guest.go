package middleware

import (
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model/entity"
	"suask/utility/resp"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

// NonGuestRequired 中间件：拒绝 guest 用户访问特定接口。
// 必须在 SessionRequired 之后使用，依赖其设置的 CtxId。
func NonGuestRequired(r *ghttp.Request) {
	userId := gconv.Int(r.Context().Value(consts.CtxId))

	var user entity.Users
	err := dao.Users.Ctx(r.Context()).
		Where(dao.Users.Columns().Id, userId).
		Fields(dao.Users.Columns().Role).
		Scan(&user)
	if err != nil || user.Role == consts.GUEST {
		resp.Do(r, 403, "请升级为正式用户后使用此功能", nil)
		return
	}

	r.Middleware.Next()
}
