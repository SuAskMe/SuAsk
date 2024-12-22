package teacher

import (
	"context"
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
	if err != nil {
		return nil, err
	}
	return res, nil
}
