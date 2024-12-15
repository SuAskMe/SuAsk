package star

import (
	"context"
	v1 "suask/api/star/v1"
	"suask/internal/service"
)

type Controller struct{}

func New() *Controller {
	return &Controller{}
}

//func (c *Controller) POST(req *ghttp.Request) {
//	req.Response.Writeln("增加收藏")
//}
//
//func (c *Controller) DELETE(req *ghttp.Request) {
//	req.Response.Writeln("删除收藏")
//}
//
//func (c *Controller) PUT(req *ghttp.Request) {
//	req.Response.Writeln("修改收藏")
//}

func (c *Controller) GET(ctx context.Context, req *v1.StarReq) (res *v1.StarRes, err error) {
	//userId = gconv.Int(ctx.Value(consts.CtxId))
	userId := 1 // 测试用
	out, err := service.Star().GetStar(ctx, userId)

	res = &v1.StarRes{}
	res.StarQuestionList = out
	return res, err
}
