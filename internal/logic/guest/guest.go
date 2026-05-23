package guest

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

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
	"github.com/gogf/gf/v2/os/gtime"
)

// ErrRateLimited is returned when the guest creation rate limit is exceeded.
// The controller layer should check this error and return HTTP 429.
var ErrRateLimited = gerror.New("请求过于频繁，请稍后再试")

// CreateGuestOutput holds the result of a successful guest creation.
type CreateGuestOutput struct {
	Type  string
	Token string
	Role  string
	Id    int
}

// CreateGuest creates a temporary guest user with device and IP rate limiting.
func CreateGuest(ctx context.Context, deviceId string, clientIP string) (*CreateGuestOutput, error) {
	if isGuestRateLimitEnabled(ctx) {
		if err := checkGuestRateLimit(ctx, deviceId, clientIP); err != nil {
			return nil, err
		}
	}

	// 3. Generate unique name: susu#XXXX (4-digit random number)
	var name string
	for i := 0; i < 10; i++ {
		name = fmt.Sprintf("susu#%04d", rand.Intn(10000))
		count, err := dao.Users.Ctx(ctx).Unscoped().Where(dao.Users.Columns().Name, name).Count()
		if err != nil {
			return nil, gerror.New(consts.ErrInternal)
		}
		if count == 0 {
			break
		}
		if i == 9 {
			return nil, gerror.New("创建临时用户失败，请重试")
		}
	}

	// 4. Insert user record (role=guest, email/salt/password_hash left nil)
	userId, err := dao.Users.Ctx(ctx).InsertAndGetId(do.Users{
		Name:     name,
		Role:     consts.GUEST,
		Nickname: name,
	})
	if err != nil {
		g.Log().Error(ctx, "CreateGuest: insert user failed", err)
		return nil, gerror.New(consts.ErrInternal)
	}

	// 5. Insert guest_users record (expires_at = now + 14 days)
	expiresAt := gtime.Now().Add(14 * 24 * time.Hour)
	_, err = dao.GuestUsers.Ctx(ctx).Insert(do.GuestUsers{
		Id:        userId,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		g.Log().Error(ctx, "CreateGuest: insert guest_users failed", err)
		// Rollback: hard-delete the user we just created
		_, _ = dao.Users.Ctx(ctx).Unscoped().Where(dao.Users.Columns().Id, userId).Delete()
		return nil, gerror.New(consts.ErrInternal)
	}

	// 6. Generate JWT and store in Redis (14-day TTL)
	token, err := sjwt.GenerateToken(int(userId))
	if err != nil {
		g.Log().Error(ctx, "CreateGuest: generate token failed", err)
		return nil, gerror.New(consts.ErrInternal)
	}
	ex := sjwt.GetExpireSecond()
	err = g.Redis().SetEX(ctx, consts.RedisJWTPrefix+strconv.Itoa(int(userId)), token, ex)
	if err != nil {
		g.Log().Error(ctx, "CreateGuest: redis setex failed", err)
		return nil, gerror.New(consts.ErrInternal)
	}

	// 7. Return token/id/role
	return &CreateGuestOutput{
		Type:  consts.TokenType,
		Token: token,
		Role:  consts.GUEST,
		Id:    int(userId),
	}, nil
}

func isGuestRateLimitEnabled(ctx context.Context) bool {
	v, err := g.Cfg().Get(ctx, "guest.rate_limit_enabled")
	if err != nil || v.IsNil() {
		return true
	}
	return v.Bool()
}

func checkGuestRateLimit(ctx context.Context, deviceId string, clientIP string) error {
	ttl := int64(g.Cfg().MustGet(ctx, "guest.rate_limit_ttl_seconds", consts.GuestRateLimitTTL).Int())
	deviceLimit := int64(g.Cfg().MustGet(ctx, "guest.device_limit", consts.GuestDeviceLimit).Int())
	ipLimit := int64(g.Cfg().MustGet(ctx, "guest.ip_limit", consts.GuestIPLimit).Int())

	if deviceId != "" {
		limited, err := incrGuestRateLimit(ctx, consts.RedisGuestDevicePrefix+deviceId, ttl, deviceLimit)
		if err != nil {
			g.Log().Error(ctx, "CreateGuest: device rate limit redis error", err)
		} else if limited {
			return ErrRateLimited
		}
	}

	if clientIP != "" {
		limited, err := incrGuestRateLimit(ctx, consts.RedisGuestIPPrefix+clientIP, ttl, ipLimit)
		if err != nil {
			g.Log().Error(ctx, "CreateGuest: IP rate limit redis error", err)
		} else if limited {
			return ErrRateLimited
		}
	}
	return nil
}

func incrGuestRateLimit(ctx context.Context, key string, ttl int64, limit int64) (bool, error) {
	cnt, err := g.Redis().Incr(ctx, key)
	if err != nil {
		return false, err
	}
	if cnt == 1 {
		_, _ = g.Redis().Expire(ctx, key, ttl)
	}
	return cnt > limit, nil
}

// UpgradeInput holds the parameters for upgrading a guest to a student.
type UpgradeInput struct {
	Name     string
	Email    string
	Password string
	Code     string
}

// UpgradeOutput holds the result of a successful guest upgrade.
type UpgradeOutput struct {
	Type  string
	Token string
	Role  string
	Id    int
}

// UpgradeGuest upgrades a temporary guest user to a full student account.
func UpgradeGuest(ctx context.Context, userId int, in *UpgradeInput) (*UpgradeOutput, error) {
	// 1. Verify current user is guest
	var user entity.Users
	err := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, userId).Scan(&user)
	if err != nil || user.Role != consts.GUEST {
		return nil, gerror.New("仅临时用户可升级")
	}

	// 2. Verify email code from Redis (same mechanism as register)
	code, err := g.Redis().Get(ctx, consts.RedisSendCodePrefix+in.Email)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	storedCode := code.String()
	if storedCode == "" {
		return nil, gerror.New("验证码已过期，请重新获取")
	}
	if storedCode != in.Code {
		cnt, _ := g.Redis().Decr(ctx, consts.RedisCountCodePrefix+in.Email)
		if cnt <= 0 {
			g.Redis().Del(ctx, consts.RedisSendCodePrefix+in.Email, consts.RedisCountCodePrefix+in.Email)
		}
		return nil, gerror.New("验证码错误")
	}
	// One-time use: delete code after successful verification
	g.Redis().Del(ctx, consts.RedisSendCodePrefix+in.Email, consts.RedisCountCodePrefix+in.Email)

	// 3. Check name uniqueness (including soft-deleted users, exclude self)
	nameCount, err := dao.Users.Ctx(ctx).Unscoped().
		Where(dao.Users.Columns().Name, in.Name).
		WhereNot(dao.Users.Columns().Id, userId).
		Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if nameCount > 0 {
		return nil, gerror.New("用户名已存在")
	}

	// 4. Check email uniqueness (including soft-deleted users)
	emailCount, err := dao.Users.Ctx(ctx).Unscoped().
		Where(dao.Users.Columns().Email, in.Email).
		Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if emailCount > 0 {
		return nil, gerror.New("邮箱已存在")
	}

	// 5. Hash password with bcrypt
	hash, err := utility.HashPassword(in.Password)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 6. Update user record: role→student, set name/nickname/email/password_hash, clear salt
	_, err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, userId).Data(do.Users{
		Name:         in.Name,
		Nickname:     in.Name,
		Email:        in.Email,
		PasswordHash: hash,
		Salt:         "",
		Role:         consts.STUDENT,
	}).Update()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 7. Delete guest_users record (no longer a guest)
	_, _ = dao.GuestUsers.Ctx(ctx).Where(dao.GuestUsers.Columns().Id, userId).Delete()

	// 8. Insert settings record (theme=default, notify_switch=on, notify_email=email)
	_, err = service.Setting().AddSetting(ctx, model.AddSettingInput{
		Id:           userId,
		ThemeId:      consts.DefaultThemeId,
		NotifySwitch: true,
		NotifyEmail:  in.Email,
	})
	if err != nil {
		g.Log().Error(ctx, "UpgradeGuest: add setting failed", err)
		// Don't fail the upgrade for a settings insert error
	}

	// 9. Re-issue JWT (role changed, need new token) and update Redis
	newToken, err := sjwt.GenerateToken(userId)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	ex := sjwt.GetExpireSecond()
	_ = g.Redis().SetEX(ctx, consts.RedisJWTPrefix+strconv.Itoa(userId), newToken, ex)

	// 10. Return new token/id/role
	return &UpgradeOutput{
		Type:  consts.TokenType,
		Token: newToken,
		Role:  consts.STUDENT,
		Id:    userId,
	}, nil
}
