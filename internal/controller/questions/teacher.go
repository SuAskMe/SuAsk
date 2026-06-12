package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	qutil "suask/internal/logic/questions_util"
	"suask/internal/middleware"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/validation"

	"github.com/gogf/gf/v2/util/gconv"
)

type cTeacherQuestion struct{}

var TeacherQuestion = cTeacherQuestion{}

func GetQuestionOfTeacherImpl(ctx context.Context, req interface{}) (res interface{}, err error) {
	baseInput := model.GetBaseOfTeacherInput{}
	gconv.Scan(req, &baseInput)

	// 管理员模式：跳过提问箱权限检查
	if !middleware.IsAdminMode(ctx) {
		err = validation.TeacherPerm(ctx, baseInput.TeacherID)
		if err != nil {
			return
		}
	}

	userId := gconv.Int(ctx.Value(consts.CtxId))

	baseInput.UserId = userId
	baseOutput, err := service.TeacherQuestion().GetBase(ctx, &baseInput)
	if err != nil {
		return
	}
	QuestionList := baseOutput.Questions
	idMap := baseOutput.IdMap
	// 获取图片 (#2 优化：批量查)
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	if imagesOutput != nil && len(imagesOutput.ImageMap) > 0 {
		allFileIDs := qutil.CollectFileIDs(imagesOutput.ImageMap, nil)
		urlMap, err_ := qutil.BatchGetFileURLs(ctx, allFileIDs)
		if err_ != nil {
			return nil, err_
		}
		imageURLs := qutil.ResolveImageURLs(urlMap, imagesOutput.ImageMap)
		for qid, urls := range imageURLs {
			QuestionList[idMap[qid]].ImageURLs = urls
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
	res_, err := GetQuestionOfTeacherImpl(ctx, req)
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
	res_, err := GetQuestionOfTeacherImpl(ctx, req)
	if err != nil {
		return
	}
	gconv.Scan(res_, &res)
	return
}
