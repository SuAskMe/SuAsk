package cmd

import (
	"context"
	"suask/internal/controller/favorite"
	"suask/internal/controller/file"
	"suask/internal/controller/history"
	"suask/internal/controller/notification"
	"suask/internal/controller/questions"
	"suask/internal/controller/register"
	"suask/internal/controller/teacher"
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
					service.Middleware().CORS,
				)
				// 这里无需登录，不需要请求用户数据
				group.Bind(
					register.Register,
					user.User.GetUserInfoById,
					teacher.Teacher.GetTeacher,
					questions.TeacherQuestion,
					notification.Notification,
				)
				// 这里是登录和非登录共有接口
				group.Group("/", func(group *ghttp.RouterGroup) {
					err := Middleware(gfToken, ctx, group)
					if err != nil {
						panic(err)
					}
					group.Bind(
						file.File.GetFileById,
						file.File.GetFileList,
						user.User.Info,
						questions.PublicQuestions.Get,
						questions.PublicQuestions.Add,
						questions.PublicQuestions.GetKeywords,
						questions.PublicQuestions.GetByKeyword,
						questions.QuestionDetail.GetDetail,
						user.User.UpdateUserInfo,
						user.User.UpdatePassWord,
						file.File.UpdateFile,
						questions.QuestionDetail.AddAnswer,
						favorite.Favorite,
						history.History,
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
