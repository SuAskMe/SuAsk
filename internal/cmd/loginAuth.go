package cmd

import (
	"context"
	"strconv"
	"strings"
	v1 "suask/api/login/v1"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model/entity"
	"suask/utility/login"
	"suask/utility/response"

	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

func LoginToken() (gfToken *gtoken.GfToken, err error) {
	gfToken = &gtoken.GfToken{
		CacheMode:        consts.CacheMode,
		ServerName:       consts.ServerName,
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFuncFrontend,
		LoginAfterFunc:   loginAfterFunc,
		LogoutPath:       "/user/logout",
		MultiLogin:       true,
		AuthAfterFunc:    authAfterFunc,
		AuthPaths:        g.SliceStr{},
		AuthExcludePaths: g.SliceStr{},
	}
	err = gfToken.Start()
	return
}

//func MustLoginToken() (gfToken *gtoken.GfToken, err error) {
//	gfToken = &gtoken.GfToken{
//		CacheMode:        consts.CacheMode,
//		ServerName:       consts.ServerName,
//		LoginPath:        "/login",
//		LoginBeforeFunc:  loginFuncFrontend,
//		LoginAfterFunc:   loginAfterFunc,
//		LogoutPath:       "/user/logout",
//		MultiLogin:       true,
//		AuthAfterFunc:    authAfterFuncMustLogin,
//		AuthPaths:        g.SliceStr{},
//		AuthExcludePaths: g.SliceStr{},
//	}
//	err = gfToken.Start()
//	return
//}

func loginFuncFrontend(r *ghttp.Request) (string, interface{}) {
	name := r.Get("name").String()
	email := r.Get("email").String()
	password := r.Get("password").String()
	ctx := context.TODO()
	// 输入参数为空
	if (password == "" && name == "") || (password == "" && email == "") {
		r.Response.WriteJson(gtoken.Fail(consts.ErrLoginFaulMsg))
		r.ExitAll()
	}
	userInfo := entity.Users{}
	var err error
	if name != "" {
		err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Name, name).Scan(&userInfo)
	}
	if email != "" {
		err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Email, email).Scan(&userInfo)
	}
	// 查不到用户
	if err != nil {
		//r.Response.WriteJson(gtoken.Fail(err.Error()))
		r.Response.WriteJson(response.JsonRes{Code: -1, Message: consts.ErrLoginFaulMsg, Data: nil})
		r.ExitAll()
	}
	// 密码校验失败
	if login.EncryptPassword(password, userInfo.Salt) != userInfo.PasswordHash {
		r.Response.WriteJson(gtoken.Fail(consts.ErrLoginFaulMsg))
		r.ExitAll()
	}
	return strconv.Itoa(userInfo.Id), userInfo
}

func loginAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	if !respData.Success() {
		respData.Code = 0
		r.Response.WriteJson(respData)
		return
	} else {
		respData.Code = 1
		userId := respData.GetString("userKey")
		userInfo := entity.Users{}
		err := dao.Users.Ctx(context.TODO()).WherePri(userId).
			Fields(dao.Users.Columns().Id).
			Fields(dao.Users.Columns().Role).
			Scan(&userInfo)
		if err != nil {
			return
		}
		data := &v1.LoginRes{
			Token: respData.GetString("token"),
			Type:  consts.TokenType,
		}
		data.Id = userInfo.Id
		data.Role = userInfo.Role
		//data.Name = userInfo.Name
		//data.Email = userInfo.Email
		//data.Nickname = userInfo.Nickname
		//data.Introduction = userInfo.Introduction
		//data.AvatarURL = avatarURL
		//data.ThemeId = userInfo.ThemeId
		response.JsonExit(r, 0, "login success", data)
	}
	return
}

func authAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	var userInfo v1.LoginRes
	err := gconv.Struct(respData.GetString("data"), &userInfo)
	//fmt.Println("Not Must Login", respData)
	if err != nil {
		if isMustLoginPath(r.URL.String()) {
			//fmt.Println("login fail, block")
			response.Auth(r)
		}
		//fmt.Println("login fail, set 1")
		r.SetCtxVar(consts.CtxId, consts.DefaultUserId)
	} else {
		//fmt.Println("login success")
		r.SetCtxVar(consts.CtxId, userInfo.Id)
	}

	// fmt.Println("login_id", r.GetCtxVar(consts.CtxId))

	r.Middleware.Next()
}

func isMustLoginPath(path string) bool {
	// 这里添加必须登陆的路由，前缀匹配
	mustLoginPath := []string{
		"/files",
		"/user/info",
		"/favorite",
		"/history",
		"/notification",
		"/questions/public",
		"/teacher/question",
	}
	for _, v := range mustLoginPath {
		if strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}
