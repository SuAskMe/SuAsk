package login

import (
	"context"
	"strconv"
	v1 "suask/api/login/v1"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/module/sjwt"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
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
	res = &v1.LogoutRes{UserId: userId}
	return
}

func (c *cLogin) HeartBeats(ctx context.Context, req *v1.HeartBeatsReq) (res *v1.HeartBeatsRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	res = &v1.HeartBeatsRes{UserId: userId}

	key := consts.RedisJWTPrefix + strconv.Itoa(userId)
	maxTTL := sjwt.GetExpireSecond() // 14 天（秒）
	threshold := maxTTL / 2          // 7 天（秒）

	// 检查 Redis TTL 剩余时间
	ttl, err := g.Redis().TTL(ctx, key)
	if err != nil || ttl <= 0 {
		err = nil // 不向调用方暴露内部错误，key 不存在或已过期时不续期
		return
	}

	// 仅当剩余 ≤ 7 天时才刷新为 14 天
	if ttl <= threshold {
		_, _ = g.Redis().Expire(ctx, key, maxTTL)

		// Guest: also renew guest_users.expires_at
		var user entity.Users
		_ = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, userId).Fields(dao.Users.Columns().Role).Scan(&user)
		if user.Role == consts.GUEST {
			newExpire := gtime.Now().Add(14 * 24 * time.Hour)
			_, _ = dao.GuestUsers.Ctx(ctx).Where(dao.GuestUsers.Columns().Id, userId).Data(g.Map{"expires_at": newExpire}).Update()
		}
	}

	return
}

// func (c *cLogin) RefreshToken(ctx context.Context, req *v1.RefreshTokenReq) (res *v1.RefreshTokenRes, err error) {
// 	return
// }
