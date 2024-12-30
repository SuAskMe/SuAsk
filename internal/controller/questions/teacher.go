package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cTeacherQuestion struct{}

var TeacherQuestion = cTeacherQuestion{}

func GetQuestionOfTeacherImpl(ctx context.Context, req interface{}) (res interface{}, err error) {
	baseInput := model.GetBaseOfTeacherInput{}
	gconv.Scan(req, &baseInput)
	baseOutput, err := service.TeacherQuestion().GetBase(ctx, &baseInput)
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
	if imagesOutput != nil {
		for k, v := range imagesOutput.ImageMap {
			urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
			if err_ != nil {
				return nil, err_
			}
			QuestionList[idMap[k]].ImageURLs = urls.URL
		}
	}
	// 返回结果
	res = &v1.GetPageOfTeacherRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
	}
	return
}

func (cTeacherQuestion) Get(ctx context.Context, req *v1.GetPageOfTeacherReq) (res *v1.GetPageOfTeacherRes, err error) {
	res_, err := GetQuestionImpl(ctx, req)
	if err != nil {
		return
	}
	gconv.Scan(res_, &res)
	return
}

func (cTeacherQuestion) GetKeywords(ctx context.Context, req *v1.GetSearchKeywordsOfTeacherReq) (res *v1.GetSearchKeywordsOfTeacherRes, err error) {
	input := model.GetKeywordsOfTeacherInput{}
	gconv.Scan(req, &input)
	ouput, err := service.TeacherQuestion().GetKeyword(ctx, &input)
	gconv.Scan(ouput, &res)
	return
}

func (cTeacherQuestion) GetByKeyword(ctx context.Context, req *v1.GetPageByKeywordOfTeacherReq) (res *v1.GetPageByKeywordOfTeacherRes, err error) {
	res_, err := GetQuestionImpl(ctx, req)
	if err != nil {
		return
	}
	gconv.Scan(res_, &res)
	return
}

//func (cTeacherQuestion) Favorite(ctx context.Context, req *v1.FavoriteOfTeacherReq) (res *v1.FavoriteOfTeacherRes, err error) {
//	input := model.FavoriteInput{}
//	gconv.Scan(req, &input)
//	output, err := service.QuestionUtil().Favorite(ctx, &input)
//	gconv.Scan(output, &res)
//	return
//}
