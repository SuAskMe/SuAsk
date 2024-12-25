package teacher

import (
	"context"
	v1 "suask/api/teacher/v1"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/entity"
	"suask/internal/service"
)

type sTeacher struct {
}

func (s *sTeacher) GetTeacher(ctx context.Context, _ model.TeacherGetInput) (out model.TeacherGetOutput, err error) {
	var teacherList []entity.Teachers
	var count int
	err = dao.Teachers.Ctx(ctx).ScanAndCount(&teacherList, &count, false)
	if err != nil {
		return model.TeacherGetOutput{}, err
	}
	out = model.TeacherGetOutput{
		TeacherList: make([]v1.TeacherBase, count),
	}
	for i, teacher := range teacherList {
		out.TeacherList[i] = v1.TeacherBase{
			Id:           teacher.Id,
			Responses:    teacher.Responses,
			Name:         teacher.Name,
			AvatarUrl:    teacher.AvatarUrl,
			Introduction: teacher.Introduction,
			Email:        teacher.Email,
			Perm:         teacher.Perm,
		}
	}
	return out, nil
}

func init() {
	service.RegisterTeacher(New())
}

func New() *sTeacher {
	return &sTeacher{}
}
