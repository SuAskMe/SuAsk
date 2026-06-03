package login

import (
	"context"
	"strconv"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/module/sjwt"
	"suask/utility"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sLogin struct{}

func (s sLogin) Login(ctx context.Context, in *model.UserLoginInput) (res *model.UserLoginOutput, err error) {
	// 输入参数为空
	if in.Name == "" && in.Email == "" {
		return nil, gerror.New("用户名和邮箱不能同时为空")
	}
	userInfo := entity.Users{}
	md := dao.Users.Ctx(ctx)
	if in.Name != "" {
		md = md.Where(dao.Users.Columns().Name, in.Name)
	} else if in.Email != "" {
		md = md.Where(dao.Users.Columns().Email, in.Email)
	}
	err = md.Scan(&userInfo)
	// 查不到用户
	if err != nil {
		return nil, gerror.New("登录失败，用户名或密码错误")
	}
	// 密码校验：兼容老 MD5 + 新 bcrypt 两种存储
	match, needUpgrade, verifyErr := utility.VerifyPassword(userInfo.PasswordHash, userInfo.Salt, in.Password)
	if verifyErr != nil {
		g.Log().Error(ctx, verifyErr)
		return nil, gerror.New(consts.ErrInternal)
	}
	if !match {
		return nil, gerror.New("登录失败，用户名或密码错误")
	}
	// 旧哈希透明升级到 bcrypt：下次登录就完全走 bcrypt 路径
	if needUpgrade {
		if newHash, hashErr := utility.HashPassword(in.Password); hashErr == nil {
			_, updErr := dao.Users.Ctx(ctx).
				Where(dao.Users.Columns().Id, userInfo.Id).
				Update(do.Users{PasswordHash: newHash, Salt: ""})
			if updErr != nil {
				// 升级失败不影响登录，仅记录，下次登录再试
				g.Log().Warningf(ctx, "bcrypt upgrade failed for user %d: %v", userInfo.Id, updErr)
			}
		} else {
			g.Log().Warningf(ctx, "bcrypt hash gen failed for user %d: %v", userInfo.Id, hashErr)
		}
	}
	// 查看是否已登录
	vtoken, err := g.Redis().Get(ctx, consts.RedisJWTPrefix+strconv.Itoa(userInfo.Id))
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	token := vtoken.String()
	// 已在别处登录，返回同一个token
	if token != "" {
		return &model.UserLoginOutput{Type: consts.TokenType, Id: userInfo.Id, Role: userInfo.Role, Token: token}, nil
	}
	// 生成token
	token, err = sjwt.GenerateToken(userInfo.Id)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("登录失败，生成token失败")
	}
	ex := sjwt.GetExpireSecond()
	err = g.Redis().SetEX(ctx, consts.RedisJWTPrefix+strconv.Itoa(userInfo.Id), token, ex)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	return &model.UserLoginOutput{Type: consts.TokenType, Id: userInfo.Id, Role: userInfo.Role, Token: token}, nil
}

func (s sLogin) Logout(ctx context.Context) error {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == 0 {
		return gerror.New("未找到登录用户")
	}
	_, err := g.Redis().Del(ctx, consts.RedisJWTPrefix+strconv.Itoa(userId))
	if err != nil {
		g.Log().Error(ctx, err)
		return gerror.New(consts.ErrInternal)
	}
	return nil
}

func (s sLogin) HeartBeats(ctx context.Context) error {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == 0 {
		return gerror.New("未找到登录用户")
	}
	key := consts.RedisJWTPrefix + strconv.Itoa(userId)
	maxTTL := sjwt.GetExpireSecond()
	threshold := maxTTL / 2

	ttl, err := g.Redis().TTL(ctx, key)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil
	}
	if ttl <= 0 || ttl > threshold {
		return nil
	}
	_, _ = g.Redis().Expire(ctx, key, maxTTL)

	var user entity.Users
	_ = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, userId).Fields(dao.Users.Columns().Role).Scan(&user)
	if user.Role == consts.GUEST {
		newExpire := gtime.Now().Add(14 * 24 * time.Hour)
		_, _ = dao.GuestUsers.Ctx(ctx).
			Where(dao.GuestUsers.Columns().Id, userId).
			Data(do.GuestUsers{ExpiresAt: newExpire}).
			Update()
	}
	return nil
}

func init() {
	service.RegisterLogin(New())
}

func New() *sLogin {
	return &sLogin{}
}
