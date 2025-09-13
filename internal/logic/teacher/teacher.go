package teacher

import (
	"context"
	v1 "suask/api/teacher/v1"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
)

type sTeacher struct {
}

func (s *sTeacher) GetTeacherList(ctx context.Context, _ model.TeacherGetInput) (out model.TeacherGetOutput, err error) {
	var teacherList []entity.Teachers
	var count int
	md := dao.Teachers.Ctx(ctx)
	md = md.Order("CASE perm WHEN 'public' THEN 1 WHEN 'protected' THEN 2 WHEN 'private' THEN 3 ELSE 4 END;")
	err = md.ScanAndCount(&teacherList, &count, false)
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

func (s *sTeacher) TeacherExist(ctx context.Context, TeacherId int) (name string, err error) {
	md := dao.Teachers.Ctx(ctx).WherePri(TeacherId).Fields(dao.Teachers.Columns().Name)
	err = md.Scan(&name)
	if err != nil {
		return "", err
	}
	return name, nil
}

func (s *sTeacher) UpdatePerm(ctx context.Context, in model.TeacherUpdatePermInput) (out model.TeacherUpdatePermOutput, err error) {
	update := do.Teachers{
		Perm: in.Perm,
	}
	_, err = dao.Teachers.Ctx(ctx).WherePri(in.TeacherId).Update(update)
	if err != nil {
		return model.TeacherUpdatePermOutput{}, err
	}
	out = model.TeacherUpdatePermOutput{
		TeacherId: in.TeacherId,
	}
	return out, nil
}

func init() {
	service.RegisterTeacher(New())
}

func New() *sTeacher {
	return &sTeacher{}
}
