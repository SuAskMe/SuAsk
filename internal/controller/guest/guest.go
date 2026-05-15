package guest

import (
	"context"

	v1 "suask/api/guest/v1"
	"suask/internal/consts"
	guestLogic "suask/internal/logic/guest"
	"suask/internal/middleware"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/send_email"
	"suask/utility/resp"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type cGuest struct{}

var Guest = cGuest{}

// Login creates a temporary guest user.
func (c *cGuest) Login(ctx context.Context, req *v1.GuestLoginReq) (res *v1.GuestLoginRes, err error) {
	r := g.RequestFromCtx(ctx)
	deviceId := r.Header.Get("X-Device-Id")
	clientIP := middleware.GetClientIP(r)

	out, err := guestLogic.CreateGuest(ctx, deviceId, clientIP)
	if err != nil {
		// Rate limit error → return 429
		if err == guestLogic.ErrRateLimited {
			resp.Do(r, 429, err.Error(), nil)
			return nil, nil
		}
		return nil, err
	}

	res = &v1.GuestLoginRes{
		Type:  out.Type,
		Token: out.Token,
		Role:  out.Role,
		Id:    out.Id,
	}
	return
}

// Upgrade upgrades a guest user to a full student account.
func (c *cGuest) Upgrade(ctx context.Context, req *v1.GuestUpgradeReq) (res *v1.GuestUpgradeRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))

	out, err := guestLogic.UpgradeGuest(ctx, userId, &guestLogic.UpgradeInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Code:     req.Code,
	})
	if err != nil {
		return nil, err
	}

	res = &v1.GuestUpgradeRes{
		Type:  out.Type,
		Token: out.Token,
		Role:  out.Role,
		Id:    out.Id,
	}
	return
}

// SendCode sends a verification code for guest upgrade.
// Reuses the existing email verification logic (check name/email uniqueness + send code).
func (c *cGuest) SendCode(ctx context.Context, req *v1.GuestSendCodeReq) (res *v1.GuestSendCodeRes, err error) {
	// Check if code already sent recently
	v, err := g.Redis().Get(ctx, consts.RedisSendCodePrefix+req.Email)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	if v.String() != "" {
		return &v1.GuestSendCodeRes{Msg: "验证码已发送，请注意查收或稍后再试"}, nil
	}

	// Check name and email uniqueness
	data := model.CheckEmailAndNameInput{Email: req.Email, Name: req.Name}
	out, err := service.Register().CheckEmailAndName(ctx, data)
	if err != nil {
		return nil, err
	}
	if out.NameDuplicated || out.EmailDuplicated {
		return &v1.GuestSendCodeRes{Msg: "邮箱或用户名重复"}, nil
	}

	// Send verification code
	code, err := send_email.SendCode(req.Email)
	if err != nil {
		return nil, err
	}
	err = g.Redis().SetEX(ctx, consts.RedisSendCodePrefix+req.Email, code, 300)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}
	err = g.Redis().SetEX(ctx, consts.RedisCountCodePrefix+req.Email, 10, 300)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New(consts.ErrInternal)
	}

	return &v1.GuestSendCodeRes{Msg: "200"}, nil
}
