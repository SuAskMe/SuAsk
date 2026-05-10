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

			// 显式把请求体上限设到 32 MiB。
			// 之前只改 config.yaml 的 server.clientMaxBodySize 会出现"1KB 上传也触发 413"，
			// 推测是配置反序列化路径或部署时配置没生效；这里用代码兜底一遍，任何配置状态都能保证上限正确。
			// 同时把 upload.max_bytes 也读出来做 sanity check。
			const minBodySize int64 = 32 * 1024 * 1024
			bodyLimit := minBodySize
			if v, err := g.Cfg().Get(ctx, "upload.max_bytes"); err == nil && !v.IsNil() {
				if n := v.Int64(); n > bodyLimit {
					bodyLimit = n
				}
			}
			s.SetClientMaxBodySize(bodyLimit)
			g.Log().Infof(ctx, "HTTP body size limit set to %d bytes", bodyLimit)

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
