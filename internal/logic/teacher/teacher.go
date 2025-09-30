package teacher

import (
	"context"
	v1 "suask/api/teacher/v1"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sTeacher struct {
}

func (s *sTeacher) GetTeacherList(ctx context.Context, _ model.TeacherGetInput) (out *model.TeacherGetOutput, err error) {
	var teacherList []v1.TeacherBase
	err = dao.Teachers.Ctx(ctx).Scan(&teacherList)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("获取教师列表失败")
	}
	out = &model.TeacherGetOutput{
		TeacherList: teacherList,
	}
	return out, nil
}

func (s *sTeacher) GetTeacherAvatar(ctx context.Context, in *model.TeacherGetAvatarInput) (out *model.TeacherGetAvatarOutput, err error) {
	err = dao.Teachers.Ctx(ctx).Where(dao.Teachers.Columns().Id, in.TeacherId).
		Fields(dao.Teachers.Columns().AvatarUrl).Scan(&out)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("获取教师头像失败")
	}
	return
}

func (s *sTeacher) TeacherExist(ctx context.Context, TeacherId int) (name string, err error) {
	md := dao.Teachers.Ctx(ctx).Where(dao.Teachers.Columns().Id, TeacherId).Fields(dao.Teachers.Columns().Name)
	var teacher entity.Teachers
	err = md.Scan(&teacher)
	if err != nil {
		g.Log().Debug(ctx, "teacher not exist", TeacherId)
		return "", err
	}
	return teacher.Name, nil
}

func (s *sTeacher) UpdatePerm(ctx context.Context, in model.TeacherUpdatePermInput) (out *model.TeacherUpdatePermOutput, err error) {
	update := do.Teachers{
		Perm: in.Perm,
	}
	_, err = dao.Teachers.Ctx(ctx).Where(dao.Teachers.Columns().Id, in.TeacherId).Update(update)
	if err != nil {
		return nil, err
	}
	out = &model.TeacherUpdatePermOutput{
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
