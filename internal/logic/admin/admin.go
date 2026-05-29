package admin

import (
	"context"
	"strconv"
	v1 "suask/api/admin/v1"
	"suask/internal/consts"
	"suask/internal/dao"
	fileLogic "suask/internal/logic/file"
	qutil "suask/internal/logic/questions_util"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/utility"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

const pageSize = 20

const (
	adminDeletedStatusAll       = "all"
	adminDeletedStatusDeleted   = "deleted"
	adminDeletedStatusUndeleted = "undeleted"
)

const adminQuestionFields = `
	q.id,
	q.title,
	q.contents,
	q.src_user_id,
	src.name AS src_user_name,
	src.nickname AS src_user_nickname,
	q.dst_user_id,
	dst.name AS dst_user_name,
	dst.nickname AS dst_user_nickname,
	q.created_at,
	q.deleted_at,
	q.views,
	q.reply_cnt,
	(
		SELECT COUNT(1)
		FROM answers a
		WHERE a.question_id = q.id
		  AND a.deleted_at IS NULL
	) AS answer_count`

type adminQuestionRow struct {
	Id              int         `orm:"id"`
	Title           string      `orm:"title"`
	Contents        string      `orm:"contents"`
	SrcUserId       int         `orm:"src_user_id"`
	SrcUserName     string      `orm:"src_user_name"`
	SrcUserNickname string      `orm:"src_user_nickname"`
	DstUserId       int         `orm:"dst_user_id"`
	DstUserName     string      `orm:"dst_user_name"`
	DstUserNickname string      `orm:"dst_user_nickname"`
	CreatedAt       *gtime.Time `orm:"created_at"`
	DeletedAt       *gtime.Time `orm:"deleted_at"`
	Views           int         `orm:"views"`
	ReplyCnt        int         `orm:"reply_cnt"`
	AnswerCount     int         `orm:"answer_count"`
}

type adminAnswerRow struct {
	Id           int         `orm:"id"`
	QuestionId   int         `orm:"question_id"`
	UserId       int         `orm:"user_id"`
	UserName     string      `orm:"user_name"`
	UserNickname string      `orm:"user_nickname"`
	UserRole     string      `orm:"user_role"`
	Contents     string      `orm:"contents"`
	CreatedAt    *gtime.Time `orm:"created_at"`
	Upvotes      int         `orm:"upvotes"`
	InReplyTo    int         `orm:"in_reply_to"`
	DeletedAt    *gtime.Time `orm:"deleted_at"`
}

type matchedAnswerCountRow struct {
	QuestionId int `orm:"question_id"`
	Count      int `orm:"cnt"`
}

// ListUsers 分页查询用户列表，支持角色筛选和关键词搜索，排除软删除用户
func ListUsers(ctx context.Context, page int, role string, keyword string) (res *v1.ListUsersRes, err error) {
	m := dao.Users.Ctx(ctx)

	// 角色筛选
	if role != "" {
		m = m.Where(dao.Users.Columns().Role, role)
	}

	// 关键词搜索（name / nickname / email 模糊匹配）
	if keyword != "" {
		likePattern := "%" + keyword + "%"
		m = m.Where(
			m.Builder().WhereOrLike(dao.Users.Columns().Name, likePattern).
				WhereOrLike(dao.Users.Columns().Nickname, likePattern).
				WhereOrLike(dao.Users.Columns().Email, likePattern),
		)
	}

	// 获取总数
	total, err := m.Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 计算剩余页数
	totalPages := (total + pageSize - 1) / pageSize
	remainPage := totalPages - page
	if remainPage < 0 {
		remainPage = 0
	}

	// 分页查询
	var users []entity.Users
	err = m.Fields(
		dao.Users.Columns().Id,
		dao.Users.Columns().Name,
		dao.Users.Columns().Nickname,
		dao.Users.Columns().Email,
		dao.Users.Columns().Role,
		dao.Users.Columns().Introduction,
		dao.Users.Columns().AvatarFileId,
		dao.Users.Columns().CreatedAt,
	).
		OrderDesc(dao.Users.Columns().CreatedAt).
		Page(page, pageSize).
		Scan(&users)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 批量解析头像 URL
	var avatarFileIDs []int
	for _, u := range users {
		if u.AvatarFileId != 0 {
			avatarFileIDs = append(avatarFileIDs, u.AvatarFileId)
		}
	}
	avatarURLMap, _ := qutil.BatchGetFileURLs(ctx, avatarFileIDs)

	// 构造返回列表
	list := make([]v1.AdminUserItem, 0, len(users))
	for _, u := range users {
		item := v1.AdminUserItem{
			Id:           u.Id,
			Name:         u.Name,
			Nickname:     u.Nickname,
			Email:        u.Email,
			Role:         u.Role,
			Introduction: u.Introduction,
		}
		if u.AvatarFileId != 0 {
			if url, ok := avatarURLMap[u.AvatarFileId]; ok {
				item.AvatarURL = url
			}
		}
		if item.AvatarURL == "" {
			item.AvatarURL = consts.DefaultAvatarURL
		}
		if u.CreatedAt != nil {
			item.CreatedAt = u.CreatedAt.String()
		}
		list = append(list, item)
	}

	res = &v1.ListUsersRes{
		List:       list,
		Total:      total,
		RemainPage: remainPage,
	}
	return
}

// CreateUser 创建用户，处理唯一性检查、密码哈希、教师表联动
func CreateUser(ctx context.Context, req *v1.CreateUserReq) (res *v1.CreateUserRes, err error) {
	// 检查用户名唯一性（包含软删除用户）
	nameCount, err := dao.Users.Ctx(ctx).Unscoped().
		Where(dao.Users.Columns().Name, req.Name).
		Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if nameCount > 0 {
		return nil, gerror.New("用户名已存在")
	}

	// 检查邮箱唯一性（包含软删除用户）
	emailCount, err := dao.Users.Ctx(ctx).Unscoped().
		Where(dao.Users.Columns().Email, req.Email).
		Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if emailCount > 0 {
		return nil, gerror.New("邮箱已存在")
	}

	// 密码哈希（bcrypt，salt 留空）
	hash, err := utility.HashPassword(req.Password)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 插入用户
	userId, err := dao.Users.Ctx(ctx).InsertAndGetId(do.Users{
		Name:         req.Name,
		Email:        req.Email,
		Salt:         "",
		PasswordHash: hash,
		Role:         req.Role,
		Nickname:     req.Nickname,
		Introduction: req.Introduction,
	})
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 如果角色是 teacher，创建 teachers 表记录
	if req.Role == consts.TEACHER {
		perm := req.Perm
		if perm == "" {
			perm = consts.PermPublic
		}
		_, err = dao.Teachers.Ctx(ctx).Insert(do.Teachers{
			Id:           userId,
			Perm:         perm,
			Responses:    0,
			Introduction: req.Introduction,
		})
		if err != nil {
			g.Log().Error(ctx, "CreateUser: insert teacher failed", "userId", userId, "err", err)
		}
	}

	res = &v1.CreateUserRes{Id: int(userId)}
	return
}

// UpdateUser 编辑用户信息，处理角色变更时的教师表联动
func UpdateUser(ctx context.Context, req *v1.UpdateUserReq, currentUserId int) (res *v1.UpdateUserRes, err error) {
	// 查询目标用户
	var targetUser entity.Users
	err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, req.Id).
		Scan(&targetUser)
	if err != nil || targetUser.Id == 0 {
		return nil, gerror.New("用户不存在")
	}

	// 禁止自我降级
	if req.Id == currentUserId && req.Role != "" && req.Role != consts.ADMIN {
		return nil, gerror.New("不允许降级自己的角色")
	}

	// 邮箱唯一性检查（如果修改了邮箱，包含软删除用户）
	if req.Email != "" && req.Email != targetUser.Email {
		emailCount, err := dao.Users.Ctx(ctx).Unscoped().
			Where(dao.Users.Columns().Email, req.Email).
			WhereNot(dao.Users.Columns().Id, req.Id).
			Count()
		if err != nil {
			return nil, gerror.New(consts.ErrInternal)
		}
		if emailCount > 0 {
			return nil, gerror.New("邮箱已被其他用户使用")
		}
	}

	// 构造更新数据
	updateData := do.Users{}
	hasUpdate := false
	if req.Nickname != "" {
		updateData.Nickname = req.Nickname
		hasUpdate = true
	}
	if req.Email != "" {
		updateData.Email = req.Email
		hasUpdate = true
	}
	if req.Role != "" {
		updateData.Role = req.Role
		hasUpdate = true
	}
	if req.Introduction != "" {
		updateData.Introduction = req.Introduction
		hasUpdate = true
	}

	// 更新用户表
	if hasUpdate {
		_, err = dao.Users.Ctx(ctx).
			Where(dao.Users.Columns().Id, req.Id).
			Data(updateData).
			Update()
		if err != nil {
			return nil, gerror.New(consts.ErrInternal)
		}
	}

	// 处理角色变更时的教师表联动
	if req.Role != "" && req.Role != targetUser.Role {
		// 原来是 teacher，现在不是 → 删除 teachers 记录
		if targetUser.Role == consts.TEACHER && req.Role != consts.TEACHER {
			_, err = dao.Teachers.Ctx(ctx).
				Where(dao.Teachers.Columns().Id, req.Id).
				Delete()
			if err != nil {
				g.Log().Error(ctx, "UpdateUser: delete teacher failed", "userId", req.Id, "err", err)
			}
		}
		// 原来不是 teacher，现在是 → 插入 teachers 记录
		if targetUser.Role != consts.TEACHER && req.Role == consts.TEACHER {
			perm := req.Perm
			if perm == "" {
				perm = consts.PermPublic
			}
			_, err = dao.Teachers.Ctx(ctx).Insert(do.Teachers{
				Id:           req.Id,
				Perm:         perm,
				Responses:    0,
				Introduction: req.Introduction,
			})
			if err != nil {
				g.Log().Error(ctx, "UpdateUser: insert teacher failed", "userId", req.Id, "err", err)
			}
		}
	}

	// 如果用户已经是 teacher 且没有角色变更，更新 teacher 相关字段
	if (req.Role == "" || req.Role == consts.TEACHER) && targetUser.Role == consts.TEACHER {
		teacherUpdate := do.Teachers{}
		hasTeacherUpdate := false
		if req.Introduction != "" {
			teacherUpdate.Introduction = req.Introduction
			hasTeacherUpdate = true
		}
		if req.Perm != "" {
			teacherUpdate.Perm = req.Perm
			hasTeacherUpdate = true
		}
		if hasTeacherUpdate {
			_, err = dao.Teachers.Ctx(ctx).
				Where(dao.Teachers.Columns().Id, req.Id).
				Data(teacherUpdate).
				Update()
			if err != nil {
				g.Log().Error(ctx, "UpdateUser: update teacher fields failed", "userId", req.Id, "err", err)
			}
		}
	}

	res = &v1.UpdateUserRes{Id: req.Id}
	return
}

// ResetPassword 重置用户密码，生成新 hash，清除 Redis JWT
func ResetPassword(ctx context.Context, userId int, newPassword string) (res *v1.ResetPasswordRes, err error) {
	// 检查用户是否存在且未软删除
	count, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userId).
		Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if count == 0 {
		return nil, gerror.New("用户不存在")
	}

	// 生成新密码哈希（bcrypt，salt 留空）
	hash, err := utility.HashPassword(newPassword)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 更新密码
	_, err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userId).
		Data(do.Users{
			Salt:         "",
			PasswordHash: hash,
		}).
		Update()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 清除 Redis JWT，强制用户重新登录
	_, err = g.Redis().Del(ctx, consts.RedisJWTPrefix+strconv.Itoa(userId))
	if err != nil {
		g.Log().Error(ctx, "ResetPassword: del redis jwt failed", "userId", userId, "err", err)
	}

	res = &v1.ResetPasswordRes{Id: userId}
	return
}

// DeleteUser 软删除用户，禁止自删，清除 Redis JWT
func DeleteUser(ctx context.Context, userId int, currentUserId int) (res *v1.DeleteUserRes, err error) {
	// 禁止自删
	if userId == currentUserId {
		return nil, gerror.New("不允许删除自己")
	}

	// 检查用户是否存在且未软删除
	count, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userId).
		Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if count == 0 {
		return nil, gerror.New("用户不存在")
	}

	// 软删除：GoFrame 自动设置 deleted_at
	_, err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userId).
		Delete()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 清除 Redis JWT，强制用户下线
	_, err = g.Redis().Del(ctx, consts.RedisJWTPrefix+strconv.Itoa(userId))
	if err != nil {
		g.Log().Error(ctx, "DeleteUser: del redis jwt failed", "userId", userId, "err", err)
	}

	res = &v1.DeleteUserRes{Id: userId}
	return
}

// UpdateAvatar 管理员修改用户头像
func UpdateAvatar(ctx context.Context, userId int) (res *v1.UpdateAvatarRes, err error) {
	// 检查用户是否存在
	count, err := dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, userId).
		Count()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if count == 0 {
		return nil, gerror.New("用户不存在")
	}

	// 从请求中获取上传的文件
	r := g.RequestFromCtx(ctx)
	avatarFile := r.GetUploadFile("avatar")

	// 调用共享的头像上传逻辑
	out, err := fileLogic.UploadAvatar(ctx, fileLogic.UploadAvatarInput{
		UserId: userId,
		File:   avatarFile,
	})
	if err != nil {
		return nil, err
	}

	res = &v1.UpdateAvatarRes{Id: out.UserId, AvatarURL: out.AvatarURL}
	return
}

// ListQuestions 分页查询管理员内容管理的一层问题列表。
func ListQuestions(ctx context.Context, req *v1.ListQuestionsReq) (res *v1.ListQuestionsRes, err error) {
	md := adminQuestionModel(ctx)
	if req.TeacherId > 0 {
		md = md.Where("q.dst_user_id = ?", req.TeacherId)
	}

	switch req.Status {
	case "", "all":
		md = md.Where("q.deleted_at IS NULL")
	case "answered":
		md = md.Where("q.deleted_at IS NULL")
		md = md.Where("EXISTS (SELECT 1 FROM answers ax WHERE ax.question_id = q.id AND ax.deleted_at IS NULL)")
	case "unanswered":
		md = md.Where("q.deleted_at IS NULL")
		md = md.Where("NOT EXISTS (SELECT 1 FROM answers ax WHERE ax.question_id = q.id AND ax.deleted_at IS NULL)")
	case "deleted":
		md = md.Where("q.deleted_at IS NOT NULL")
	default:
		return nil, gerror.New("无效的问题状态")
	}

	likePattern := "%" + req.Keyword + "%"
	if req.Keyword != "" {
		md = md.Where(`
			(
				q.title LIKE ?
				OR q.contents LIKE ?
				OR src.name LIKE ?
				OR src.nickname LIKE ?
				OR dst.name LIKE ?
				OR dst.nickname LIKE ?
				OR EXISTS (
				SELECT 1
				FROM answers ax
				WHERE ax.question_id = q.id
				  AND ax.deleted_at IS NULL
				  AND ax.contents LIKE ?
				)
			)`,
			likePattern, likePattern, likePattern, likePattern, likePattern, likePattern, likePattern,
		)
	}

	var rows []adminQuestionRow
	var total int
	err = md.Order("q.created_at DESC").
		Page(req.Page, pageSize).
		ScanAndCount(&rows, &total, false)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	matchedAnswerCounts := map[int]int{}
	if req.Keyword != "" && len(rows) > 0 {
		qIDs := make([]int, 0, len(rows))
		for _, row := range rows {
			qIDs = append(qIDs, row.Id)
		}
		var counts []matchedAnswerCountRow
		err = g.DB().Ctx(ctx).Model("answers").
			Fields("question_id, COUNT(1) AS cnt").
			WhereIn("question_id", qIDs).
			Where("deleted_at IS NULL").
			WhereLike("contents", likePattern).
			Group("question_id").
			Scan(&counts)
		if err != nil {
			return nil, gerror.New(consts.ErrInternal)
		}
		for _, count := range counts {
			matchedAnswerCounts[count.QuestionId] = count.Count
		}
	}

	list := make([]v1.AdminQuestionItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, buildAdminQuestionItem(row, matchedAnswerCounts[row.Id]))
	}

	totalPages := (total + pageSize - 1) / pageSize
	remainPage := totalPages - req.Page
	if remainPage < 0 {
		remainPage = 0
	}

	return &v1.ListQuestionsRes{
		List:       list,
		Total:      total,
		RemainPage: remainPage,
	}, nil
}

// GetQuestionDetail 返回管理员视角的问题详情，回答作为二级内容挂在问题下。
func GetQuestionDetail(ctx context.Context, req *v1.GetQuestionDetailReq) (res *v1.GetQuestionDetailRes, err error) {
	deletedStatus, err := resolveAdminDeletedStatus(req.DeletedStatus, req.IncludeDeleted)
	if err != nil {
		return nil, err
	}

	var question adminQuestionRow
	err = adminQuestionModel(ctx).
		Where("q.id = ?", req.Id).
		Scan(&question)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if question.Id == 0 {
		return nil, gerror.New("问题不存在")
	}

	md := g.DB().Ctx(ctx).Model("answers a").
		Unscoped().
		LeftJoin("users u", "u.id = a.user_id AND u.deleted_at IS NULL").
		Fields(`
			a.id,
			a.question_id,
			a.user_id,
			u.name AS user_name,
			u.nickname AS user_nickname,
			u.role AS user_role,
			a.contents,
			a.created_at,
			a.upvotes,
			a.in_reply_to,
			a.deleted_at`).
		Where("a.question_id = ?", req.Id)
	md = applyDeletedStatusFilter(md, "a.deleted_at", deletedStatus)

	var answerRows []adminAnswerRow
	err = md.Order("a.created_at ASC").Scan(&answerRows)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	answers := make([]v1.AdminQuestionAnswerItem, 0, len(answerRows))
	for _, row := range answerRows {
		answers = append(answers, buildAdminQuestionAnswerItem(row))
	}

	return &v1.GetQuestionDetailRes{
		Question: buildAdminQuestionItem(question, 0),
		Answers:  answers,
	}, nil
}

// DeleteQuestion 管理员软删除问题，复用已有的问题删除权限与行为。
func DeleteQuestion(ctx context.Context, questionId int, currentUserId int) (res *v1.DeleteQuestionRes, err error) {
	if err = service.QuestionDetail().DeleteQuestion(ctx, questionId, currentUserId); err != nil {
		return nil, err
	}
	return &v1.DeleteQuestionRes{Id: questionId}, nil
}

// RestoreQuestion 管理员恢复已删除问题。
func RestoreQuestion(ctx context.Context, questionId int) (res *v1.RestoreQuestionRes, err error) {
	var question entity.Questions
	err = dao.Questions.Ctx(ctx).
		Unscoped().
		Fields(dao.Questions.Columns().Id, dao.Questions.Columns().DeletedAt).
		Where(dao.Questions.Columns().Id, questionId).
		Scan(&question)
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	if question.Id == 0 {
		return nil, gerror.New("问题不存在")
	}
	if question.DeletedAt == nil {
		return nil, gerror.New("问题未删除")
	}

	_, err = dao.Questions.Ctx(ctx).
		Unscoped().
		Where(dao.Questions.Columns().Id, questionId).
		Data(dao.Questions.Columns().DeletedAt, gdb.Raw("NULL")).
		Update()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}
	return &v1.RestoreQuestionRes{Id: questionId}, nil
}

// DeleteQuestionAnswer 管理员软删除问题下的回答，并同步维护问题当前可见回复数。
func DeleteQuestionAnswer(ctx context.Context, questionId int, answerId int) (res *v1.DeleteQuestionAnswerRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var answer entity.Answers
		if err := tx.Model("answers").
			Fields("id, question_id").
			Where("id = ? AND question_id = ? AND deleted_at IS NULL", answerId, questionId).
			Scan(&answer); err != nil {
			return gerror.New(consts.ErrInternal)
		}
		if answer.Id == 0 {
			return gerror.New("回答不存在")
		}

		if _, err := tx.Exec(
			"UPDATE answers SET deleted_at = CURRENT_TIMESTAMP WHERE id = ? AND question_id = ? AND deleted_at IS NULL",
			answerId,
			questionId,
		); err != nil {
			return gerror.New(consts.ErrInternal)
		}

		if _, err := tx.Exec(
			"UPDATE questions SET reply_cnt = CASE WHEN reply_cnt > 0 THEN reply_cnt - 1 ELSE 0 END WHERE id = ?",
			questionId,
		); err != nil {
			return gerror.New(consts.ErrInternal)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.DeleteQuestionAnswerRes{Id: answerId, QuestionId: questionId}, nil
}

// RestoreQuestionAnswer 管理员恢复问题下的已删除回答，并同步维护问题当前可见回复数。
func RestoreQuestionAnswer(ctx context.Context, questionId int, answerId int) (res *v1.RestoreQuestionAnswerRes, err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		var answer struct {
			Id         int         `orm:"id"`
			QuestionId int         `orm:"question_id"`
			DeletedAt  *gtime.Time `orm:"deleted_at"`
		}
		if err := tx.Model("answers").
			Unscoped().
			Fields("id, question_id, deleted_at").
			Where("id = ? AND question_id = ?", answerId, questionId).
			Scan(&answer); err != nil {
			return gerror.New(consts.ErrInternal)
		}
		if answer.Id == 0 {
			return gerror.New("回答不存在")
		}
		if answer.DeletedAt == nil {
			return gerror.New("回答未删除")
		}

		if _, err := tx.Exec(
			"UPDATE answers SET deleted_at = NULL WHERE id = ? AND question_id = ? AND deleted_at IS NOT NULL",
			answerId,
			questionId,
		); err != nil {
			return gerror.New(consts.ErrInternal)
		}

		if _, err := tx.Exec(
			"UPDATE questions SET reply_cnt = reply_cnt + 1 WHERE id = ?",
			questionId,
		); err != nil {
			return gerror.New(consts.ErrInternal)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &v1.RestoreQuestionAnswerRes{Id: answerId, QuestionId: questionId}, nil
}

func adminQuestionModel(ctx context.Context) *gdb.Model {
	return g.DB().Ctx(ctx).Model("questions q").
		Unscoped().
		LeftJoin("users src", "src.id = q.src_user_id AND src.deleted_at IS NULL").
		LeftJoin("users dst", "dst.id = q.dst_user_id AND dst.deleted_at IS NULL").
		Fields(adminQuestionFields)
}

func buildAdminQuestionItem(row adminQuestionRow, matchedAnswerCount int) v1.AdminQuestionItem {
	status := "unanswered"
	if row.AnswerCount > 0 {
		status = "answered"
	}
	srcName := row.SrcUserName
	srcNickname := row.SrcUserNickname
	if srcName == "" && srcNickname == "" {
		srcNickname = deletedUserLabel(row.SrcUserId)
	}
	dstName := row.DstUserName
	dstNickname := row.DstUserNickname
	if dstName == "" && dstNickname == "" {
		dstNickname = deletedTeacherLabel(row.DstUserId)
	}
	return v1.AdminQuestionItem{
		Id:                 row.Id,
		Title:              row.Title,
		Contents:           row.Contents,
		SrcUserId:          row.SrcUserId,
		SrcUserName:        srcName,
		SrcUserNickname:    srcNickname,
		DstUserId:          row.DstUserId,
		DstUserName:        dstName,
		DstUserNickname:    dstNickname,
		CreatedAt:          timeMilli(row.CreatedAt),
		Views:              row.Views,
		ReplyCnt:           row.ReplyCnt,
		AnswerCount:        row.AnswerCount,
		MatchedAnswerCount: matchedAnswerCount,
		Status:             status,
		IsDeleted:          row.DeletedAt != nil,
		DeletedAt:          timeMilli(row.DeletedAt),
	}
}

func buildAdminQuestionAnswerItem(row adminAnswerRow) v1.AdminQuestionAnswerItem {
	userName := row.UserName
	userNickname := row.UserNickname
	if userName == "" && userNickname == "" {
		userNickname = deletedUserLabel(row.UserId)
	}
	return v1.AdminQuestionAnswerItem{
		Id:           row.Id,
		QuestionId:   row.QuestionId,
		UserId:       row.UserId,
		UserName:     userName,
		UserNickname: userNickname,
		UserRole:     row.UserRole,
		Contents:     row.Contents,
		CreatedAt:    timeMilli(row.CreatedAt),
		Upvotes:      row.Upvotes,
		InReplyTo:    row.InReplyTo,
		IsDeleted:    row.DeletedAt != nil,
		DeletedAt:    timeMilli(row.DeletedAt),
	}
}

func deletedUserLabel(userId int) string {
	return "未知用户"
}

func deletedTeacherLabel(userId int) string {
	return "未知用户"
}

func timeMilli(t *gtime.Time) int64 {
	if t == nil {
		return 0
	}
	return t.TimestampMilli()
}

func resolveAdminDeletedStatus(deletedStatus string, includeDeleted bool) (string, error) {
	switch deletedStatus {
	case "":
		if includeDeleted {
			return adminDeletedStatusAll, nil
		}
		return adminDeletedStatusUndeleted, nil
	case adminDeletedStatusAll, adminDeletedStatusDeleted, adminDeletedStatusUndeleted:
		return deletedStatus, nil
	default:
		return "", gerror.New("无效的删除状态")
	}
}

func applyDeletedStatusFilter(md *gdb.Model, deletedAtField string, deletedStatus string) *gdb.Model {
	switch deletedStatus {
	case adminDeletedStatusDeleted:
		return md.Where(deletedAtField + " IS NOT NULL")
	case adminDeletedStatusUndeleted:
		return md.Where(deletedAtField + " IS NULL")
	default:
		return md
	}
}
