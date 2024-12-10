package cmd

import (
	"context"
	"suask/internal/controller/hello"
	"suask/internal/controller/register"
	"suask/internal/controller/user"
	"suask/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()

			gfToken, err := LoginToken()
			if err != nil {
				return err
			}

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					ghttp.MiddlewareHandlerResponse,
					service.Middleware().CORS)
				group.Bind(
					register.Register.Register,
					user.User.GetUserInfoById,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					err := gfToken.Middleware(ctx, group)
					if err != nil {
						panic(err)
					}
					group.Bind(
						hello.NewV1(),
						user.User.Info,
						user.User.UpdateUserInfo,
						user.User.UpdatePassWord,
					)
				})
			})
			s.Run()
			return nil
		},
	}
)
