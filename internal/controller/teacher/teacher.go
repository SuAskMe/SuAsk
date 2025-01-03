package teacher

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "suask/api/teacher/v1"
	"suask/internal/model"
	"suask/internal/service"
)

type cTeacher struct {
}

var Teacher cTeacher

func (c *cTeacher) GetTeacher(ctx context.Context, req *v1.TeacherReq) (res *v1.TeacherRes, err error) {
	out := model.TeacherGetOutput{}
	out, err = service.Teacher().GetTeacher(ctx, model.TeacherGetInput{})
	if err != nil {
		return nil, err
	}
	res = &v1.TeacherRes{}
	res.TeacherList = out.TeacherList
	return res, nil
}

func (c *cTeacher) GetTeacherPin(ctx context.Context, req *v1.TeacherPinReq) (res *v1.TeacherPinRes, err error) {
	out, err := service.TeacherQuestionSelf().GetQFMPinned(ctx, &model.GetQFMInput{TeacherId: req.TeacherId})
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
	res = &v1.TeacherPinRes{}
	err = gconv.Scan(qfm, &res.QuestionList)
	if err != nil {
		return nil, err
	}
	return res, err
}
