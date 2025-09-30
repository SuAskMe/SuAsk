package cmd

import (
	"suask/internal/middleware"
)

// var trieMux *triemux.TrieMux

// func init() {
// 	// 这里添加必须登陆的路由，前缀匹配
// 	mustLoginPath := []string{
// 		"/files",
// 		"/user/info",
// 		"/favorite",
// 		"/history",
// 		"/notification",
// 		"/questions/public",
// 		"/teacher",
// 		"/user/heartbeat",
// 	}
// 	trieMux = triemux.NewTrieMux()
// 	for _, v := range mustLoginPath {
// 		if err := trieMux.Insert(v); err != nil {
// 			panic(err)
// 		}
// 	}
// }

func JwtToken() (jm *middleware.JWTMiddleware) {
	prefixUrlMatch := []string{
		"/files",
		"/favorite",
		"/history",
		"/notification",
		"/questions/public",
		"/teacher",
		"/user",
	}
	excpetUrlMatch := []string{"/user/send-code", "/user/forget-password"}
	jm.BuildMustLoginTrie(prefixUrlMatch, nil, excpetUrlMatch)
	return jm
}

// func LoginToken() (gfToken *gtoken.GfToken, err error) {
// 	gfToken = &gtoken.GfToken{
// 		CacheMode:        1,
// 		ServerName:       consts.ServerName,
// 		LoginPath:        "/login",
// 		LoginBeforeFunc:  loginFuncFrontend,
// 		LoginAfterFunc:   loginAfterFunc,
// 		LogoutPath:       "/user/logout",
// 		MultiLogin:       true,
// 		AuthAfterFunc:    authAfterFunc,
// 		AuthPaths:        g.SliceStr{},
// 		AuthExcludePaths: g.SliceStr{},
// 	}
// 	err = gfToken.Start()
// 	return
// }

// func loginFuncFrontend(r *ghttp.Request) (string, any) {
// 	name := r.Get("name").String()
// 	email := r.Get("email").String()
// 	password := r.Get("password").String()
// 	ctx := context.TODO()
// 	// 输入参数为空
// 	if password == "" || (name == "" && email == "") {
// 		r.Response.WriteJson(gtoken.Fail("登录失败，用户名或密码错误"))
// 		r.ExitAll()
// 	}
// 	userInfo := entity.Users{}
// 	var err error
// 	if name != "" {
// 		err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Name, name).Scan(&userInfo)
// 	} else if email != "" {
// 		err = dao.Users.Ctx(ctx).Where(dao.Users.Columns().Email, email).Scan(&userInfo)
// 	}
// 	// 查不到用户
// 	if err != nil {
// 		//r.Response.WriteJson(gtoken.Fail(err.Error()))
// 		r.Response.WriteJson(response.JsonRes{Code: -1, Message: "登录失败，用户名或密码错误", Data: nil})
// 		r.ExitAll()
// 	}
// 	// 密码校验失败
// 	if utility.EncryptPassword(password, userInfo.Salt) != userInfo.PasswordHash {
// 		r.Response.WriteJson(gtoken.Fail("登录失败，用户名或密码错误"))
// 		r.ExitAll()
// 	}
// 	return strconv.Itoa(userInfo.Id), userInfo
// }

// func loginAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
// 	if !respData.Success() {
// 		respData.Code = 0
// 		r.Response.WriteJson(respData)
// 		return
// 	} else {
// 		respData.Code = 1
// 		userId := respData.GetString("userKey")
// 		userInfo := entity.Users{}
// 		err := dao.Users.Ctx(context.TODO()).Where(dao.Users.Columns().Id, userId).
// 			Fields(dao.Users.Columns().Id).
// 			Fields(dao.Users.Columns().Role).
// 			Scan(&userInfo)
// 		if err != nil {
// 			return
// 		}
// 		data := &v1.LoginRes{
// 			Token: respData.GetString("token"),
// 			Type:  consts.TokenType,
// 		}
// 		data.Id = userInfo.Id
// 		data.Role = userInfo.Role
// 		//data.Name = userInfo.Name
// 		//data.Email = userInfo.Email
// 		//data.Nickname = userInfo.Nickname
// 		//data.Introduction = userInfo.Introduction
// 		//data.AvatarURL = avatarURL
// 		//data.ThemeId = userInfo.ThemeId
// 		response.JsonExit(r, 0, "login success", data)
// 	}
// }

// func authAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
// 	var userInfo v1.LoginRes
// 	err := gconv.Struct(respData.GetString("data"), &userInfo)
// 	//fmt.Println("Not Must Login", respData)
// 	if err != nil {
// 		g.Log().Debug(context.TODO(), "authAfterFunc err", err, r.URL.String())
// 		if isMustLoginPath(r.URL.String()) {
// 			//fmt.Println("login fail, block")
// 			g.Log().Debug(context.TODO(), "is Must Login")
// 			response.Auth(r)
// 		}
// 		//fmt.Println("login fail, set 1")
// 		r.SetCtxVar(consts.CtxId, consts.DefaultUserId)
// 	} else {
// 		//fmt.Println("login success")
// 		r.SetCtxVar(consts.CtxId, userInfo.Id)
// 	}

// 	// fmt.Println("login_id", r.GetCtxVar(consts.CtxId))

// 	r.Middleware.Next()
// }

// func isMustLoginPath(path string) bool {
// 	ok := trieMux.HasPrefix(path)
// 	if ok {
// 		return ok
// 	}
// 	fullUrlMatch := []string{"/user"}
// 	return slices.Contains(fullUrlMatch, path)
// }
