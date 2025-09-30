package cmd

import (
	"context"
	"suask/internal/controller/favorite"
	"suask/internal/controller/history"
	"suask/internal/controller/login"
	"suask/internal/controller/notification"
	"suask/internal/controller/questions"
	"suask/internal/controller/register"
	"suask/internal/controller/teacher"
	"suask/internal/controller/user"
	"suask/internal/middleware"

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

			jToken := JwtToken()
			if err != nil {
				return err
			}

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					ghttp.MiddlewareHandlerResponse,
					middleware.CORS,
				)
				// 这里无需登录，不需要请求用户数据
				group.Bind(
					register.Register,
					user.User.GetUserInfoById,
					teacher.Teacher.GetTeacher,
					teacher.Teacher.GetTeacherPin,
				)
				// 这里是登录和非登录共有接口
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(jToken.JwtAuth)
					group.Bind(login.Login.Login,
						login.Login.Logout,
						login.Login.HeartBeats,
						questions.PublicQuestions,
						questions.QuestionDetail.GetDetail,
						user.User.Info,
						user.User.UpdateUserInfo,
						user.User.UpdatePassWord,
						user.User.SendVerificationCode,
						user.User.ForgetPassword,
						questions.QuestionDetail.AddAnswer,
						favorite.Favorite,
						history.History,
						questions.QuestionDetail.Upvote,
						questions.Question,
						teacher.Teacher.UpdatePerm,
						questions.TeacherSelf,
						questions.TeacherQuestion,
						notification.Notification,
					)
				})
			})
			// 设置静态文件服务
			s.SetIndexFolder(true)
			s.SetFileServerEnabled(true)
			s.SetServerRoot(".")

			// 启动服务器
			s.Run()
			return nil
		},
	}
)
