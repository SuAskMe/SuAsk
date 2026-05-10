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
			// 静态文件服务：只暴露上传目录（从配置读取 upload.path），
			// 明确拒绝把项目根目录当作静态根 —— 否则 database/*.db、
			// manifest/config/*.yaml、logs/*、main.exe 都会通过 HTTP 被下载。
			uploadPath := g.Cfg().MustGet(ctx, "upload.path").String()
			if uploadPath == "" {
				uploadPath = "upload"
			}
			// AddStaticPath 会自动开启 fileServer，不需要再调 SetServerRoot。
			s.AddStaticPath("/"+uploadPath, "./"+uploadPath)
			// 显式禁止目录列表，避免 /upload/ 被直接翻目录。
			s.SetIndexFolder(false)

			// 启动服务器
			s.Run()
			return nil
		},
	}
)
