package validation

import (
	"context"
	"errors"
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
		md := dao.Teachers.Ctx(ctx).Where("id = ?", teacherId).Fields(dao.Teachers.Columns().Perm)
		var teacher *entity.Teachers
		err := md.Scan(&teacher)
		if err != nil {
			return err
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

// 所有问题细节查看权限（不检查老师提问箱权限）
func QuestionPerm(ctx context.Context, question *entity.Questions) error {
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	if question.IsPrivate && question.SrcUserId != UserId { // 私有问题，且不是自己提问
		return errors.New("你不能查看别人的私有问题")
	}
	if question.DstUserId == 0 && UserId == consts.DefaultUserId { // 问大家的问题
		return errors.New("请登录后再查看问大家的问题")
	}
	if question.DstUserId != 0 && question.ReplyCnt <= 0 { // 问教师的问题,且还没有回复
		switch UserId {
		case consts.DefaultUserId:
			return errors.New("该问题还没有回复，请耐心等待")
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
	if question.DstUserId == 0 { // 问大家的问题
		return nil
	}
	if question.DstUserId != 0 { // 问教师的问题,且还没有回复
		switch UserId {
		case question.SrcUserId:
			if question.ReplyCnt <= 0 { // 还没有回复
				return errors.New("该问题还没有回复，请耐心等待")
			}
			return nil
		case question.DstUserId:
			return nil
		default:
			return errors.New("你不能回答这个问题")
		}
	}
	return nil
}

// 判断是否为老师
func IsTeacher(ctx context.Context, teacherId int) error {
	_, ok := teacherCache.Load(teacherId)
	if !ok {
		// fmt.Println("not in cache", teacherId)
		md := dao.Teachers.Ctx(ctx).Where("id = ?", teacherId).Fields(dao.Teachers.Columns().Perm)
		var teacher *entity.Teachers
		err := md.Scan(&teacher)
		if err != nil {
			return err
		}
		teacherCache.Store(teacherId, teacher)
		// fmt.Println(teacher)
		if teacher.Perm == "" {
			return errors.New("该用户不是老师")
		}
	}
	return nil
}

func UpdateTeacherPerm(teacherId int, perm string) {
	teacher := &entity.Teachers{Perm: perm}
	teacherCache.Store(teacherId, teacher)
}
