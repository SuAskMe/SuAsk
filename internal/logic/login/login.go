package login

import (
	"context"
	"suask/internal/model"
	"suask/internal/service"
)

type sLogin struct{}

func (s sLogin) Login(ctx context.Context, in model.UserLoginInput) error {
	//TODO implement me
	panic("implement me")
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
