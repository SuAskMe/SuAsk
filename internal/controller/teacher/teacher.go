package teacher

import (
	"context"
	v1 "suask/api/teacher/v1"
	"suask/internal/consts"
	qutil "suask/internal/logic/questions_util"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/validation"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/gconv"
)

type cTeacher struct {
}

var Teacher cTeacher

func (c *cTeacher) GetTeacher(ctx context.Context, req *v1.TeacherReq) (res *v1.TeacherRes, err error) {
	out, err := service.Teacher().GetTeacherList(ctx, model.TeacherGetInput{})
	if err != nil {
		return nil, err
	}
	res = &v1.TeacherRes{TeacherList: out.TeacherList}
	return res, nil
}

func (c *cTeacher) GetTeacherPin(ctx context.Context, req *v1.TeacherPinReq) (res *v1.TeacherPinRes, err error) {
	out, err := service.TeacherQuestionSelf().GetQFMPinned(ctx, &model.GetQFMInput{TeacherId: req.TeacherId})
	if err != nil {
		return
	}
	qfm := out.Questions
	idMap := out.IdMap
	// 获取图片（批量查）
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: out.QuestionIDs})
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
			qfm[idMap[qid]].ImageURLs = urls
		}
	}
	res = &v1.TeacherPinRes{}
	err = gconv.Scan(qfm, &res.QuestionList)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *cTeacher) UpdatePerm(ctx context.Context, req *v1.UpdatePermReq) (res *v1.UpdatePermRes, err error) {
	teacherId := gconv.Int(ctx.Value(consts.CtxId))
	name, err := service.Teacher().TeacherExist(ctx, teacherId)
	if err != nil {
		return nil, err
	}
	if len(name) == 0 {
		return nil, gerror.New("教师不存在")
	}
	_, err = service.Teacher().UpdatePerm(ctx, model.TeacherUpdatePermInput{
		TeacherId: teacherId,
		Perm:      gconv.String(req.Perm),
	})
	if err != nil {
		return nil, err
	}
	res = &v1.UpdatePermRes{
		Id: teacherId,
	}
	validation.UpdateTeacherPerm(teacherId, name, gconv.String(req.Perm))
	return res, nil
}
