package login

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
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
	// 旧哈希透明升级到 bcrypt
	if needUpgrade {
		if newHash, hashErr := utility.HashPassword(in.Password); hashErr == nil {
			_, updErr := dao.Users.Ctx(ctx).
				Where(dao.Users.Columns().Id, userInfo.Id).
				Update(do.Users{PasswordHash: newHash, Salt: ""})
			if updErr != nil {
				g.Log().Warningf(ctx, "bcrypt upgrade failed for user %d: %v", userInfo.Id, updErr)
			}
		} else {
			g.Log().Warningf(ctx, "bcrypt hash gen failed for user %d: %v", userInfo.Id, hashErr)
		}
	}

	return &model.UserLoginOutput{
		Id:   userInfo.Id,
		Role: userInfo.Role,
	}, nil
}

func (s sLogin) Logout(ctx context.Context) error {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == 0 {
		return gerror.New("未找到登录用户")
	}
	// Session deletion is now handled at the controller level (needs access to cookie).
	return nil
}

func (s sLogin) HeartBeats(ctx context.Context) error {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == 0 {
		return gerror.New("未找到登录用户")
	}

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
