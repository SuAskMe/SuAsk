package admin

import (
	"context"
	v1 "suask/api/admin/v1"
	"suask/internal/consts"
	"suask/internal/logic/admin"

	"github.com/gogf/gf/v2/util/gconv"
)

type cAdmin struct{}

var Admin cAdmin

// ListUsers 管理员-用户列表
func (c *cAdmin) ListUsers(ctx context.Context, req *v1.ListUsersReq) (res *v1.ListUsersRes, err error) {
	return admin.ListUsers(ctx, req.Page, req.Role, req.Keyword)
}

// CreateUser 管理员-创建用户
func (c *cAdmin) CreateUser(ctx context.Context, req *v1.CreateUserReq) (res *v1.CreateUserRes, err error) {
	return admin.CreateUser(ctx, req)
}

// UpdateUser 管理员-编辑用户
func (c *cAdmin) UpdateUser(ctx context.Context, req *v1.UpdateUserReq) (res *v1.UpdateUserRes, err error) {
	currentUserId := gconv.Int(ctx.Value(consts.CtxId))
	return admin.UpdateUser(ctx, req, currentUserId)
}

// ResetPassword 管理员-重置密码
func (c *cAdmin) ResetPassword(ctx context.Context, req *v1.ResetPasswordReq) (res *v1.ResetPasswordRes, err error) {
	return admin.ResetPassword(ctx, req.Id, req.Password)
}

// DeleteUser 管理员-删除用户
func (c *cAdmin) DeleteUser(ctx context.Context, req *v1.DeleteUserReq) (res *v1.DeleteUserRes, err error) {
	currentUserId := gconv.Int(ctx.Value(consts.CtxId))
	return admin.DeleteUser(ctx, req.Id, currentUserId)
}

// UpdateAvatar 管理员-修改用户头像
func (c *cAdmin) UpdateAvatar(ctx context.Context, req *v1.UpdateAvatarReq) (res *v1.UpdateAvatarRes, err error) {
	return admin.UpdateAvatar(ctx, req.Id)
}
