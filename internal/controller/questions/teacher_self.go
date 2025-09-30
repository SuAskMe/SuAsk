package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/validation"

	"github.com/gogf/gf/v2/util/gconv"
)

type cTeacherSelf struct{}

var TeacherSelf = cTeacherSelf{}

func GetQFMImpl(ctx context.Context, in *model.GetQFMInput) (res *v1.QFMBase, err error) {
	_, err = validation.IsTeacher(ctx, in.TeacherId)
	if err != nil {
		return nil, err
	}
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
	if imagesOutput != nil {
		for k, v := range imagesOutput.ImageMap {
			urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
			if err_ != nil {
				return nil, err_
			}
			qfm[idMap[k]].ImageURLs = urls.URL
		}
	}
	res = &v1.QFMBase{
		QFMList:    qfm,
		RemainPage: out.RemainPage,
	}
	return
}

func (cTeacherSelf) GetQFMAll(ctx context.Context, req *v1.GetQFMReq) (res *v1.GetQFMRes, err error) {
	Tid := gconv.Int(ctx.Value(consts.CtxId))
	// fmt.Println("GetQFMAll Tid:", Tid)
	_, err = validation.IsTeacher(ctx, Tid)
	if err != nil {
		return nil, err
	}

	// Tid := 2

	var in model.GetQFMInput
	gconv.Scan(req, &in)
	in.TeacherId = Tid
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	res = &v1.GetQFMRes{
		QFMBase: *_res,
	}
	return
}

func (cTeacherSelf) GetQFMAnswered(ctx context.Context, req *v1.GetQFMAnsweredReq) (res *v1.GetQFMAnsweredRes, err error) {
	// Tid := 2
	Tid := gconv.Int(ctx.Value(consts.CtxId))

	var in model.GetQFMInput
	gconv.Scan(req, &in)
	in.TeacherId = Tid
	in.Tag = consts.Answered
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	res = &v1.GetQFMAnsweredRes{
		QFMBase: *_res,
	}
	return
}

func (cTeacherSelf) GetQFMUnanswered(ctx context.Context, req *v1.GetQFMUnansweredReq) (res *v1.GetQFMUnansweredRes, err error) {
	// Tid := 2
	Tid := gconv.Int(ctx.Value(consts.CtxId))

	var in model.GetQFMInput
	gconv.Scan(req, &in)
	in.TeacherId = Tid
	in.Tag = consts.Unanswered
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	res = &v1.GetQFMUnansweredRes{
		QFMBase: *_res,
	}
	return
}

func (cTeacherSelf) GetQFMTop(ctx context.Context, _ *v1.GetQFMTopReq) (res *v1.GetQFMTopRes, err error) {
	Tid := gconv.Int(ctx.Value(consts.CtxId))

	_, err = validation.IsTeacher(ctx, Tid)
	if err != nil {
		return nil, err
	}

	// Tid := 2

	out, err := service.TeacherQuestionSelf().GetQFMPinned(ctx, &model.GetQFMInput{TeacherId: Tid})
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
	Tid := gconv.Int(ctx.Value(consts.CtxId))
	_, err = validation.IsTeacher(ctx, Tid)
	if err != nil {
		return nil, err
	}

	// Tid := 2

	var in model.GetQFMKeywordsInput
	gconv.Scan(req, &in)
	in.TeacherId = Tid
	out, err := service.TeacherQuestionSelf().GetKeyword(ctx, &in)
	gconv.Scan(out, &res)
	return
}

func (cTeacherSelf) GetQFMByKeyword(ctx context.Context, req *v1.SearchQFMReq) (res *v1.SearchQFMRes, err error) {
	// Tid := 2
	Tid := gconv.Int(ctx.Value(consts.CtxId))

	var in model.GetQFMInput
	gconv.Scan(req, &in)
	in.TeacherId = Tid
	_res, err := GetQFMImpl(ctx, &in)
	if err != nil {
		return
	}
	res = &v1.SearchQFMRes{
		QFMBase: *_res,
	}
	return
}

func (cTeacherSelf) PinQFMInput(ctx context.Context, req *v1.PinQFMReq) (res *v1.PinQFMRes, err error) {
	// Tid := 2
	Tid := gconv.Int(ctx.Value(consts.CtxId))

	_, err = validation.IsTeacher(ctx, Tid)
	if err != nil {
		return nil, err
	}

	out, err := service.TeacherQuestionSelf().PinQFM(ctx, &model.PinQFMInput{QuestionId: req.QuestionId, TeacherId: Tid})
	if err != nil {
		return
	}
	res = &v1.PinQFMRes{IsPinned: out.IsPinned}
	return
}
