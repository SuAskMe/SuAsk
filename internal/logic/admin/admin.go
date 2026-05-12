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
	"suask/utility"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const pageSize = 20

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
