package cmd

import (
	"context"
	"suask/internal/controller/favorite"
	"suask/internal/controller/file"
	"suask/internal/controller/hello"
	"suask/internal/controller/questions"
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
				// 这里是不需要认证的接口
				group.Bind(
					register.Register,
					user.User.GetUserInfoById,
					file.File.GetFileById,
					questions.PublicQuestions,
					favorite.Favorite.GetFavorite,
					favorite.Favorite.GetPageFavorite,
					favorite.Favorite.DelFavorite,
				)
				group.Group("/", func(group *ghttp.RouterGroup) {
					err := gfToken.Middleware(ctx, group)
					if err != nil {
						panic(err)
					}
					group.Bind(
						hello.NewV1(),
						// 这里是需要认证的接口
						user.User.Info,
						user.User.UpdateUserInfo,
						user.User.UpdatePassWord,
						file.File.UpdateFile,
					)
				})
			})
			// 设置静态文件服务
			s.SetIndexFolder(true)
			s.SetFileServerEnabled(true)
			s.SetServerRoot(".")
			//s.AddSearchPath("/Users/john/Documents")

			// 启动服务器
			s.Run()
			return nil
		},
	}
)
