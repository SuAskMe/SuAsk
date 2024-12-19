package questions

import (
	"context"
	"strconv"
	v1 "suask/api/questions/v1"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cPublicQuestions struct{}

var PublicQuestions = cPublicQuestions{}

func (cPublicQuestions) Get(ctx context.Context, req *v1.GetPageReq) (res *v1.GetPageRes, err error) {
	baseInput := model.GetBaseInput{}
	gconv.Scan(req, &baseInput)
	baseOutput, err := service.PublicQuestion().GetBase(ctx, &baseInput)
	if err != nil {
		return
	}
	QuestionList := baseOutput.Questions
	idMap := baseOutput.IdMap
	// 获取图片
	imagesOutput, err := service.PublicQuestion().GetImages(ctx, &model.GetImagesInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range imagesOutput.ImageMap {
		QuestionList[idMap[k]].ImageURLs = []string{strconv.Itoa(v[0]) + ".jpg"}
	}
	// 获取回答数
	answersOutput, err := service.PublicQuestion().GetAnswers(ctx, &model.GetAnswersInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range answersOutput.CountMap {
		QuestionList[idMap[k]].AnswerNum = v
		QuestionList[idMap[k]].AnswerAvatars = []string{strconv.Itoa(answersOutput.AvatarsMap[k][0])}
	}
	// 返回结果
	res = &v1.GetPageRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
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
	baseInput := model.GetBaseInput{}
	gconv.Scan(req, &baseInput)
	baseOutput, err := service.PublicQuestion().GetBase(ctx, &baseInput)
	if err != nil {
		return
	}
	QuestionList := baseOutput.Questions
	idMap := baseOutput.IdMap
	// 获取图片
	imagesOutput, err := service.PublicQuestion().GetImages(ctx, &model.GetImagesInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range imagesOutput.ImageMap {
		QuestionList[idMap[k]].ImageURLs = []string{strconv.Itoa(v[0]) + ".jpg"}
	}
	// 获取回答数
	answersOutput, err := service.PublicQuestion().GetAnswers(ctx, &model.GetAnswersInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range answersOutput.CountMap {
		QuestionList[idMap[k]].AnswerNum = v
		QuestionList[idMap[k]].AnswerAvatars = []string{strconv.Itoa(answersOutput.AvatarsMap[k][0])}
	}
	// 返回结果
	res = &v1.GetPageByKeywordRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
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
