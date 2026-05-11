package teacher

import (
	"context"
	v1 "suask/api/teacher/v1"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sTeacher struct {
}

func (s *sTeacher) GetTeacherList(ctx context.Context, _ model.TeacherGetInput) (out *model.TeacherGetOutput, err error) {
	// teachers 表已精简（只有 id, perm, responses, avatar_url），
	// name/email/introduction 从 users 表 JOIN 获取。
	type teacherRow struct {
		Id           int    `json:"id"`
		Responses    int    `json:"responses"`
		Name         string `json:"name"`
		AvatarUrl    string `json:"avatarUrl"`
		Introduction string `json:"introduction"`
		Email        string `json:"email"`
		Perm         string `json:"perm"`
	}
	var rows []teacherRow
	err = g.DB().Ctx(ctx).Model("teachers t").
		LeftJoin("users u", "u.id = t.id").
		Fields("t.id, t.responses, u.name, t.avatar_url, u.introduction, u.email, t.perm").
		Where("u.deleted_at IS NULL").
		Scan(&rows)
	if err != nil {
		g.Log().Error(ctx, err)
		return nil, gerror.New("获取教师列表失败")
	}
	teacherList := make([]v1.TeacherBase, len(rows))
	for i, r := range rows {
		teacherList[i] = v1.TeacherBase{
			Id:           r.Id,
			Responses:    r.Responses,
			Name:         r.Name,
			AvatarUrl:    r.AvatarUrl,
			Introduction: r.Introduction,
			Email:        r.Email,
			Perm:         r.Perm,
		}
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
	// name 现在在 users 表里
	type row struct {
		Name string `json:"name"`
	}
	var r row
	err = g.DB().Ctx(ctx).Model("users").
		Where("id = ?", TeacherId).
		Fields("name").
		Scan(&r)
	if err != nil || r.Name == "" {
		g.Log().Debug(ctx, "teacher not exist", TeacherId)
		return "", err
	}
	return r.Name, nil
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
