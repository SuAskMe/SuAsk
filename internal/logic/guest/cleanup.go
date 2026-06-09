package guest

import (
	"context"

	"suask/internal/dao"
	"suask/internal/model/entity"
	"suask/module/session"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/gtime"
)

// StartGuestCleanupTask starts a periodic task that cleans up expired guest users every hour.
func StartGuestCleanupTask(ctx context.Context) {
	_, err := gcron.Add(ctx, "0 0 * * * *", func(ctx context.Context) {
		cleanupExpiredGuests(ctx)
	}, "guest-cleanup")
	if err != nil {
		g.Log().Warning(ctx, "guest-cleanup: 定时任务注册失败", err)
	} else {
		g.Log().Info(ctx, "guest-cleanup: 定时清理任务已注册 (interval: every hour)")
	}
}

// cleanupExpiredGuests finds and hard-deletes guest users whose expires_at has passed.
func cleanupExpiredGuests(ctx context.Context) {
	// 1. Query expired guest users (expires_at < now)
	var expiredGuests []entity.GuestUsers
	err := dao.GuestUsers.Ctx(ctx).
		WhereLT(dao.GuestUsers.Columns().ExpiresAt, gtime.Now()).
		Scan(&expiredGuests)
	if err != nil {
		g.Log().Error(ctx, "guest-cleanup: 查询过期 guest 失败", err)
		return
	}
	if len(expiredGuests) == 0 {
		return
	}

	// 2. Collect user IDs
	ids := make([]int, 0, len(expiredGuests))
	for _, gu := range expiredGuests {
		ids = append(ids, gu.Id)
	}

	// 3. Hard-delete from users table (CASCADE will auto-delete guest_users rows).
	// Use raw SQL to bypass GoFrame's soft-delete behavior (users table has deleted_at).
	result, err := g.DB().Exec(ctx, "DELETE FROM users WHERE id IN (?)", ids)
	if err != nil {
		g.Log().Error(ctx, "guest-cleanup: 删除过期 guest 用户失败", err)
		return
	}
	affected, _ := result.RowsAffected()

	// 4. Clean up Redis JWT keys for deleted users
	for _, id := range ids {
		_ = session.DeleteUserSessions(ctx, id)
	}

	// 5. Log results
	if affected > 0 {
		g.Log().Infof(ctx, "guest-cleanup: 已删除 %d 个过期 guest 用户", affected)
	}
}
