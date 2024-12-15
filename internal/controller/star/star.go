package star

import (
	"context"
	v1 "suask/api/star/v1"
	"suask/internal/model"
	"suask/internal/service"
)

type Controller struct{}

var Star Controller

//func (c *Controller) POST(req *ghttp.Request) {
//	req.Response.Writeln("增加收藏")
//}

func (c *Controller) DELETE(ctx context.Context, req *v1.DeleteStarReq) (res *v1.DeleteStarRes, err error) {
	in := model.DeleteStarInput{Id: req.Id}
	out, err := service.Star().DeleteStar(ctx, in)

	res = &v1.DeleteStarRes{String: out}
	return res, err
}

//func (c *Controller) PUT(req *ghttp.Request) {
//	req.Response.Writeln("修改收藏")
//}

func (c *Controller) GET(ctx context.Context, req *v1.StarReq) (res *v1.StarRes, err error) {
	//userId = gconv.Int(ctx.Value(consts.CtxId))
	userId := 1 // 测试用
	out, err := service.Star().GetStar(ctx, userId)

	res = &v1.StarRes{StarQuestionList: out}
	return res, err
}
