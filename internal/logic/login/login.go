package login

import (
	"context"
	"strconv"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
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
	err = md.Fields(dao.Users.Columns().Id, dao.Users.Columns().Salt,
		dao.Users.Columns().PasswordHash).Scan(&userInfo)
	// 查不到用户
	if err != nil {
		return nil, gerror.New("登录失败，用户名或密码错误")
	}
	// 密码校验失败
	if utility.EncryptPassword(in.Password, userInfo.Salt) != userInfo.PasswordHash {
		return nil, gerror.New("登录失败，用户名或密码错误")
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
	//TODO implement me
	panic("implement me")
}

func (s sLogin) HeartBeats(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func init() {
	service.RegisterLogin(New())
}

func New() *sLogin {
	return &sLogin{}
}
