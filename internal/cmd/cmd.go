package cmd

import (
	"context"
	"suask/internal/controller/admin"
	"suask/internal/controller/announcement"
	"suask/internal/controller/favorite"
	"suask/internal/controller/guest"
	"suask/internal/controller/history"
	"suask/internal/controller/login"
	"suask/internal/controller/notification"
	"suask/internal/controller/questions"
	"suask/internal/controller/register"
	"suask/internal/controller/teacher"
	"suask/internal/controller/user"
	fileCleanup "suask/internal/logic/file"
	guestLogic "suask/internal/logic/guest"
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

			// 初始化校园网网段配置
			middleware.InitCampusSubnets(ctx)

			// 请求体上限 32 MiB
			const minBodySize int64 = 32 * 1024 * 1024
			bodyLimit := minBodySize
			if v, err := g.Cfg().Get(ctx, "upload.max_bytes"); err == nil && !v.IsNil() {
				if n := v.Int64(); n > bodyLimit {
					bodyLimit = n
				}
			}
			s.SetClientMaxBodySize(bodyLimit)
			g.Log().Infof(ctx, "HTTP body size limit set to %d bytes", bodyLimit)

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(
					ghttp.MiddlewareHandlerResponse,
					middleware.CORS,
				)

				// ========== 公开接口（需要校园网） ==========
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.CampusNetworkCheck)
					group.Bind(register.Register)
				})

				// ========== 公开接口（无需认证、无需校园网） ==========
				group.Bind(
					user.User.GetUserInfoById,
					teacher.Teacher.GetTeacher,
					teacher.Teacher.GetTeacherPin,
					questions.HotQuestion,
					announcement.Announcement.List,
					announcement.Announcement.Detail,
					guest.Guest.Login,
					login.Login.Login,
				)

				// ========== Guest 升级（需要校园网 + 登录） ==========
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.CampusNetworkCheck)
					group.Middleware(middleware.JwtRequired)
					group.Bind(
						guest.Guest.Upgrade,
						guest.Guest.SendCode,
					)
				})

				// ========== 必须登录（Guest 也可以访问） ==========
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.JwtRequired)
					group.Bind(
						login.Login.Logout,
						login.Login.HeartBeats,
						user.User.Info,
						questions.TeacherQuestion,
						questions.QuestionDetail.GetDetail,
						questions.Question,
						questions.QuestionDetail.DeleteAnswer,
						questions.QuestionDetail.AddAnswer,
						questions.QuestionDetail.Upvote,
						questions.Inbox,
						history.History,
						announcement.Announcement.AddComment,
					)
				})

				// ========== 必须登录 + 非 Guest（Guest 不可访问） ==========
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.JwtRequired)
					group.Middleware(middleware.NonGuestRequired)
					group.Bind(
						user.User.UpdateUserInfo,
						user.User.UpdatePassWord,
						user.User.SendVerificationCode,
						user.User.ForgetPassword,
						user.User.Deactivate,
						favorite.Favorite,
						teacher.Teacher.UpdatePerm,
						notification.Notification,
						announcement.Announcement.Create,
						announcement.Announcement.Update,
						announcement.Announcement.Delete,
					)
				})

				// ========== 管理员接口（需要登录 + 管理员角色） ==========
				group.Group("/", func(group *ghttp.RouterGroup) {
					group.Middleware(middleware.JwtRequired)
					group.Middleware(middleware.AdminRequired)
					group.Bind(
						admin.Admin.ListUsers,
						admin.Admin.CreateUser,
						admin.Admin.UpdateUser,
						admin.Admin.ResetPassword,
						admin.Admin.DeleteUser,
						admin.Admin.UpdateAvatar,
					)
				})
			})

			// 静态文件服务
			uploadPath := g.Cfg().MustGet(ctx, "upload.path").String()
			if uploadPath == "" {
				uploadPath = "upload"
			}
			s.AddStaticPath("/"+uploadPath, "./"+uploadPath)
			s.SetIndexFolder(false)

			// 启动文件清理定时任务
			fileCleanup.StartCleanupTask(ctx)

			// 启动 Guest 清理定时任务
			guestLogic.StartGuestCleanupTask(ctx)

			s.Run()
			return nil
		},
	}
)
