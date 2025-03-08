package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cPublicQuestions struct{}

// 这里的api要求严格登陆才能访问，所以这里不需要做权限验证
var PublicQuestions = cPublicQuestions{}

func GetQuestionImpl(ctx context.Context, req interface{}) (res interface{}, err error) {
	baseInput := model.GetBaseInput{}
	gconv.Scan(req, &baseInput)
	baseOutput, err := service.PublicQuestion().GetBase(ctx, &baseInput)
	if err != nil {
		return
	}
	QuestionList := baseOutput.Questions
	idMap := baseOutput.IdMap
	// 获取图片
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range imagesOutput.ImageMap {
		urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
		if err_ != nil {
			return nil, err_
		}
		QuestionList[idMap[k]].ImageURLs = urls.URL
	}
	// 获取回答数
	answersOutput, err := service.PublicQuestion().GetAnswers(ctx, &model.GetAnswersInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	if answersOutput != nil {
		for k, v := range answersOutput.AvatarsMap {
			idList := make([]int, 0, len(v))
			URLs := make([]string, 0, len(v))
			for _, u := range v {
				if u != 0 {
					idList = append(idList, u)
				} else {
					URLs = append(URLs, consts.DefaultAvatarURL)
				}
			}
			urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: idList})
			if err_ != nil {
				return nil, err_
			}
			if len(URLs) == 0 {
				URLs = urls.URL
			} else {
				URLs = append(urls.URL, URLs[0])
			}
			QuestionList[idMap[k]].AnswerAvatars = URLs
		}
	}
	// 返回结果
	res = &v1.GetPageRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
	}
	return
}

func (cPublicQuestions) Get(ctx context.Context, req *v1.GetPageReq) (res *v1.GetPageRes, err error) {
	res_, err := GetQuestionImpl(ctx, req)
	if err != nil {
		return
	}
	gconv.Scan(res_, &res)
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
	res_, err := GetQuestionImpl(ctx, req)
	if err != nil {
		return
	}
	gconv.Scan(res_, &res)
	return
}

//func (cPublicQuestions) Favorite(ctx context.Context, req *v1.FavoriteReq) (res *v1.FavoriteRes, err error) {
//	input := model.FavoriteInput{}
//	gconv.Scan(req, &input)
//	output, err := service.QuestionUtil().Favorite(ctx, &input)
//	gconv.Scan(output, &res)
//	return
//}
