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

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
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

// Logout 仅作为 service.ILogin 接口契约存在；真正的登出逻辑在
// internal/controller/login/login.go 里直接操作 Redis 完成。
// 这里以前写的是 panic("implement me")，若被误调会直接 500 + goroutine crash，
// 现在改为 no-op 并记录一条 debug 日志，方便定位谁在错误调用。
func (s sLogin) Logout(ctx context.Context) error {
	g.Log().Debug(ctx, "sLogin.Logout called; logout is handled in controller layer")
	return nil
}

// HeartBeats 同 Logout：占位实现，真逻辑在 controller 层。
func (s sLogin) HeartBeats(ctx context.Context) error {
	g.Log().Debug(ctx, "sLogin.HeartBeats called; heartbeat is handled in controller layer")
	return nil
}

func init() {
	service.RegisterLogin(New())
}

func New() *sLogin {
	return &sLogin{}
}
