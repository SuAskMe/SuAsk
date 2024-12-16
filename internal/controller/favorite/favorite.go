package favorite

import (
	"context"
	v1 "suask/api/favorite/v1"
	"suask/internal/model"
	"suask/internal/service"
)

type cFavorite struct{}

var Favorite cFavorite

//func (c *Controller) POST(req *ghttp.Request) {
//	req.Response.Writeln("增加收藏")
//}

func (c *cFavorite) DelFavorite(ctx context.Context, req *v1.DeleteFavoriteReq) (res *v1.DeleteFavoriteRes, err error) {
	in := model.DeleteFavoriteInput{Id: req.Id}
	out, err := service.Favorite().DeleteFavorite(ctx, in)

	res = &v1.DeleteFavoriteRes{String: out.String} // 从FavoriteOut中取出string，放到Res中
	return res, err
}

//func (c *Controller) PUT(req *ghttp.Request) {
//	req.Response.Writeln("修改收藏")
//}

func (c *cFavorite) GetFavorite(ctx context.Context, req *v1.FavoriteReq) (res *v1.FavoriteRes, err error) {
	//userId = gconv.Int(ctx.Value(consts.CtxId))
	userId := 1 // 测试用
	out, err := service.Favorite().GetFavorite(ctx, userId)

	res = &v1.FavoriteRes{FavoriteQuestionList: out.FavoriteQuestionList}
	return res, err
}
