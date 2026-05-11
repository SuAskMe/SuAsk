package validation

import (
	"context"
	"errors"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model/entity"
	"sync"

	"github.com/gogf/gf/v2/util/gconv"
)

var teacherCache = sync.Map{}

/*
	 问题权限：
	 1. 公开问题：
	    只要登陆了，就可以有任何权限

	 2. 问老师的问题：
	    对于查看、收藏：
			老师提问箱设置了公开
			老师回复了问题之后，所有人才能看，否则只有提问者和老师可以看(除默认用户外)

	    对于回答：（点赞只要登陆就都可以）
			只有在老师回答之后，提问者才可以回答（除默认用户外）
			其它人都不能回答

		对于提问：
			老师提问箱设置了公开，登录和未登录都可以提问

	 3. 判断是否为老师
*/

// 老师提问箱查看和提问权限
func TeacherPerm(ctx context.Context, teacherId int) error {
	if teacherId == 0 {
		return nil
	}
	t, ok := teacherCache.Load(teacherId)
	if !ok {
		md := dao.Teachers.Ctx(ctx).Where("id = ?", teacherId).Fields(dao.Teachers.Columns().Name, dao.Teachers.Columns().Perm)
		var teacher *entity.Teachers
		err := md.Scan(&teacher)
		if err != nil {
			return err
		}
		if teacher == nil || teacher.Perm == "" {
			return fmt.Errorf("该用户不是老师")
		}
		teacherCache.Store(teacherId, teacher)
		t = teacher
	}
	teacher := t.(*entity.Teachers)
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	switch teacher.Perm {
	case consts.PermPublic:
		return nil
	case consts.PermPrivate:
		return errors.New("老师并未开启提问箱，请联系老师")
	case consts.PermProtected:
		if UserId == consts.DefaultUserId {
			return errors.New("请登录后再提问")
		}
		return nil
	default: // 未知权限
		return errors.New("未知权限")
	}
}

/*
	问题权限（"问大家"模块已下线后的简化版）：
	1. 查看：
		问老师的问题 + 老师已回复 → 任何（含匿名）可看
		问老师的问题 + 老师未回复 → 只有提问者本人和该老师可看
		私有问题                  → 只有提问者本人可看
	2. 回答：
		未登录用户不能回答
		只有提问者本人在老师已回复后才能回答
		老师自己随时能回答
	3. 所有问题的 DstUserId 必定 > 0（DB 有 NOT NULL 约束）
*/

// 所有问题细节查看权限（不检查老师提问箱权限）
func QuestionPerm(ctx context.Context, question *entity.Questions) error {
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	if question.IsPrivate && question.SrcUserId != UserId { // 私有问题，且不是自己提问
		return errors.New("你不能查看别人的私有问题")
	}
	// 问老师的问题，还没有回复
	if question.ReplyCnt <= 0 {
		switch UserId {
		case question.SrcUserId:
			return nil
		case question.DstUserId:
			return nil
		default:
			return errors.New("该问题还没有回复，请耐心等待")
		}
	}
	return nil
}

// 回答问题权限 (不检查老师提问箱权限)
func AnswerPerm(ctx context.Context, question *entity.Questions) error {
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	if UserId == consts.DefaultUserId {
		return errors.New("请登录后再回答问题")
	}
	// 所有问题都是问老师的问题
	switch UserId {
	case question.SrcUserId:
		if question.ReplyCnt <= 0 { // 提问者必须等到老师先回复
			return errors.New("该问题还没有回复，请耐心等待")
		}
		return nil
	case question.DstUserId:
		return nil
	default:
		return errors.New("你不能回答这个问题")
	}
}

// 判断是否为老师
func IsTeacher(ctx context.Context, teacherId int) (string, error) {
	t, ok := teacherCache.Load(teacherId)
	if !ok {
		// fmt.Println("not in cache", teacherId)
		md := dao.Teachers.Ctx(ctx).Where("id = ?", teacherId).Fields(dao.Teachers.Columns().Perm)
		var teacher *entity.Teachers
		err := md.Scan(&teacher)
		if err != nil {
			return "", err
		}
		if teacher == nil || teacher.Perm == "" {
			return "", fmt.Errorf("该用户不是老师")
		}
		teacherCache.Store(teacherId, teacher)
		//fmt.Println(teacher)
		t = teacher
	}
	return t.(*entity.Teachers).Perm, nil
}

// UpdateTeacherPerm 同步更新缓存中的老师权限信息。
// 旧实现只 Store 了 Perm，会把已经缓存的 Name / Email 等字段清空，
// 导致后续 GetTeacherName / TeacherPerm 报"该用户不是老师"。
// 这里改为：若缓存已有对象则复制并覆盖 Perm；否则用调用方传入的 name 新建一条。
func UpdateTeacherPerm(teacherId int, name, perm string) {
	if v, ok := teacherCache.Load(teacherId); ok {
		if cached, ok := v.(*entity.Teachers); ok && cached != nil {
			updated := *cached // 浅拷贝，避免与其他 reader 共享同一指针
			updated.Perm = perm
			if name != "" {
				updated.Name = name
			}
			teacherCache.Store(teacherId, &updated)
			return
		}
	}
	teacherCache.Store(teacherId, &entity.Teachers{
		Id:   teacherId,
		Name: name,
		Perm: perm,
	})
}

func GetTeacherName(ctx context.Context, teacherId int) (string, error) {
	t, ok := teacherCache.Load(teacherId)
	if !ok {
		// fmt.Println("not in cache", teacherId)
		md := dao.Teachers.Ctx(ctx).Where("id = ?", teacherId).Fields(dao.Teachers.Columns().Name)
		var teacher *entity.Teachers
		err := md.Scan(&teacher)
		if err != nil {
			return "", err
		}
		if teacher == nil || teacher.Name == "" {
			return "", fmt.Errorf("该用户不是老师")
		}
		teacherCache.Store(teacherId, teacher)
		t = teacher
	}
	return t.(*entity.Teachers).Name, nil
}
