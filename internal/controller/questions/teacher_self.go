package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cTeacherSelf struct{}

var TeacherSelf = cTeacherSelf{}

func GetQFMImpl(ctx context.Context, in *model.GetQFMInput) (res interface{}, err error) {
	out, err := service.TeacherQuestionSelf().GetQFMAll(ctx, in)
	if err != nil {
		return
	}
	qfm := out.Questions
	idMap := out.IdMap
	// 获取图片
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: out.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range imagesOutput.ImageMap {
		urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
		if err_ != nil {
			return nil, err_
		}
		qfm[idMap[k]].ImageURLs = urls.URL
	}
	_res := &v1.GetQFMRes{}
	_res.QFMList = qfm
	_res.RemainPage = out.RemainPage
	return _res, nil
}

func (cTeacherSelf) GetQFMAll(ctx context.Context, req *v1.GetQFMReq) (res *v1.GetQFMRes, err error) {
	var in model.GetQFMInput
	gconv.Scan(req, &in)
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	gconv.Scan(_res, &res)
	return
}

func (cTeacherSelf) GetQFMAnswered(ctx context.Context, req *v1.GetQFMAnsweredReq) (res *v1.GetQFMAnsweredRes, err error) {
	var in model.GetQFMInput
	gconv.Scan(req, &in)
	in.Tag = consts.Answered
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	gconv.Scan(_res, &res)
	return
}

func (cTeacherSelf) GetQFMUnanswered(ctx context.Context, req *v1.GetQFMUnansweredReq) (res *v1.GetQFMUnansweredRes, err error) {
	var in model.GetQFMInput
	gconv.Scan(req, &in)
	in.Tag = consts.Unanswered
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	gconv.Scan(_res, &res)
	return
}

func (cTeacherSelf) GetQFMTop(ctx context.Context, _ *v1.GetQFMTopReq) (res *v1.GetQFMTopRes, err error) {
	out, err := service.TeacherQuestionSelf().GetQFMPinned(ctx, nil)
	if err != nil {
		return
	}
	qfm := out.Questions
	idMap := out.IdMap
	// 获取图片
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: out.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range imagesOutput.ImageMap {
		urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
		if err_ != nil {
			return nil, err_
		}
		qfm[idMap[k]].ImageURLs = urls.URL
	}
	res = &v1.GetQFMTopRes{}
	res.QFMList = qfm
	res.RemainPage = out.RemainPage
	return
}

func (cTeacherSelf) GetQFMKeywords(ctx context.Context, req *v1.GetQFMSearchKeywordsReq) (res *v1.GetQFMSearchKeywordsRes, err error) {
	var in model.GetQFMKeywordsInput
	gconv.Scan(req, &in)
	out, err := service.TeacherQuestionSelf().GetKeyword(ctx, &in)
	if err != nil {
		return
	}
	gconv.Scan(out, &res)
	return
}

func (cTeacherSelf) GetQFMByKeyword(ctx context.Context, req *v1.SearchQFMReq) (res *v1.SearchQFMRes, err error) {
	var in model.GetQFMInput
	gconv.Scan(req, &in)
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	gconv.Scan(_res, &res)
	return
}
