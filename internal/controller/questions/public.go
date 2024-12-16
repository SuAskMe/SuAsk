package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cPublicQuestions struct{}

var PublicQuestions = cPublicQuestions{}

func (cPublicQuestions) Get(ctx context.Context, req *v1.GetPageReq) (res *v1.GetPageRes, err error) {
	input := model.GetInput{}
	gconv.Scan(req, &input)
	// fmt.Println(input)
	ouput, err := service.PublicQuestion().Get(ctx, &input)
	// fmt.Println(err)
	res = &v1.GetPageRes{
		QuestionList: ouput.Questions,
		RemainPage:   ouput.RemainPage,
	}
	return
}

func (cPublicQuestions) GetKeywords(ctx context.Context, req *v1.GetSearchKeywordsReq) (res *v1.GetSearchKeywordsRes, err error) {
	input := model.GetKeywordsInput{}
	gconv.Scan(req, &input)
	ouput, err := service.PublicQuestion().GetKeyword(ctx, &input)
	gconv.Scan(ouput, &res)
	return
}

func (cPublicQuestions) GetByKeyword(ctx context.Context, req *v1.GetPageByKeywordReq) (res *v1.GetPageByKeywordRes, err error) {
	input := model.GetInput{}
	gconv.Scan(req, &input)
	ouput, err := service.PublicQuestion().Get(ctx, &input)
	res = &v1.GetPageByKeywordRes{
		QuestionList: ouput.Questions,
		RemainPage:   ouput.RemainPage,
	}
	return
}

func (cPublicQuestions) Favorite(ctx context.Context, req *v1.FavoriteReq) (res *v1.FavoriteRes, err error) {
	input := model.FavoriteInput{}
	gconv.Scan(req, &input)
	output, err := service.PublicQuestion().Favorite(ctx, &input)
	gconv.Scan(output, &res)
	return
}
