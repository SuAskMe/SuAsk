package star

import "github.com/gogf/gf/v2/net/ghttp"

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) POST(req *ghttp.Request) {
	req.Response.Writeln("增加收藏")
}

func (c *Controller) DELETE(req *ghttp.Request) {
	req.Response.Writeln("删除收藏")
}

func (c *Controller) PUT(req *ghttp.Request) {
	req.Response.Writeln("修改收藏")
}

func (c *Controller) GET(req *ghttp.Request) {
	req.Response.Writeln("查询收藏")
}
