package response

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
)

type JsonRes struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Json(r *ghttp.Request, code int, message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	} else {
		responseData = g.Map{}
	}
	r.Response.WriteJson(JsonRes{
		Code:    code,
		Message: message,
		Data:    responseData,
	})
}

func JsonExit(r *ghttp.Request, code int, message string, data ...interface{}) {
	Json(r, code, message, data...)
	r.Exit()
}

func dataReturn(r *ghttp.Request, code int, req ...interface{}) *JsonRes {
	var msg string
	var data interface{}
	if len(req) > 0 {
		msg = gconv.String(req[0])
	}
	if len(req) > 1 {
		data = req[1]
	}
	if code != 1 && !gconv.Bool(r.GetCtxVar("api_code")) {
		code = 401
	}
	response := &JsonRes{
		Code:    code,
		Message: msg,
		Data:    data,
	}
	r.SetParam("apiReturnRes", response)
	return response
}

func Auth(r *ghttp.Request) {
	res := dataReturn(r, 999, "请登录")
	r.Response.WriteJsonExit(res)
}
