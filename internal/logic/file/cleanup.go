package file

import (
	"context"
	"os"
	"path/filepath"
	"suask/internal/dao"
	"suask/internal/model/entity"
	files "suask/utility/files"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/text/gstr"
)

// StartCleanupTask 启动定时清理无引用文件的后台任务。
// cron 表达式和文件保护时间从配置文件读取：
//   - cleanup.cron: cron 表达式，默认 "0 0 3 * * *"（每天凌晨3点）
//   - cleanup.protect_minutes: 文件创建后多少分钟内不清理，默认 60
func StartCleanupTask(ctx context.Context) {
	cronExpr := g.Cfg().MustGet(ctx, "cleanup.cron", "0 0 3 * * *").String()
	_, err := gcron.Add(ctx, cronExpr, func(ctx context.Context) {
		CleanupOrphanFiles(ctx)
	}, "file-cleanup")
	if err != nil {
		g.Log().Warning(ctx, "file-cleanup: 定时任务注册失败", err)
	} else {
		g.Log().Infof(ctx, "file-cleanup: 定时清理任务已注册，cron=%s", cronExpr)
	}
}

// CleanupOrphanFiles 清理无引用的文件。
// 查找 files 表中未被 users.avatar_file_id 或 attachments.file_id 引用的记录，
// 删除磁盘文件并软删除数据库记录。
// 只清理创建时间超过 1 小时的文件（避免删除刚上传还未关联的文件）。
func CleanupOrphanFiles(ctx context.Context) {
	g.Log().Info(ctx, "file-cleanup: 开始清理无引用文件")

	// 获取所有被引用的 file_id
	referencedIDs := getReferencedFileIDs(ctx)

	// 获取所有未软删除的文件记录（创建超过保护时间）
	var allFiles []entity.Files
	protectMinutes := g.Cfg().MustGet(ctx, "cleanup.protect_minutes", 60).Int()
	cutoff := time.Now().Add(-time.Duration(protectMinutes) * time.Minute)
	err := dao.Files.Ctx(ctx).
		Where("created_at < ?", cutoff).
		Scan(&allFiles)
	if err != nil {
		g.Log().Error(ctx, "file-cleanup: 查询文件列表失败", err)
		return
	}

	// 找出孤儿文件
	var orphanCount, deletedDisk, deletedDB int
	for _, f := range allFiles {
		if _, ok := referencedIDs[f.Id]; ok {
			continue
		}
		orphanCount++

		// 删除磁盘文件
		diskPath := getFileDiskPath(ctx, f.Hash, f.Name)
		if diskPath != "" {
			if err := os.Remove(diskPath); err == nil {
				deletedDisk++
				// 清理空目录
				cleanEmptyDirs(diskPath)
			}
		}

		// 软删除数据库记录
		_, err := dao.Files.Ctx(ctx).
			Where(dao.Files.Columns().Id, f.Id).
			Delete()
		if err == nil {
			deletedDB++
		}
	}

	g.Log().Infof(ctx, "file-cleanup: 完成。孤儿文件=%d, 删除磁盘=%d, 软删除DB=%d",
		orphanCount, deletedDisk, deletedDB)
}

// getReferencedFileIDs 获取所有被引用的 file_id 集合
func getReferencedFileIDs(ctx context.Context) map[int]struct{} {
	referenced := make(map[int]struct{})

	// users.avatar_file_id
	type avatarRow struct {
		AvatarFileId int `orm:"avatar_file_id"`
	}
	var avatars []avatarRow
	g.DB().Ctx(ctx).Model("users").
		Fields("avatar_file_id").
		Where("avatar_file_id IS NOT NULL").
		Scan(&avatars)
	for _, a := range avatars {
		referenced[a.AvatarFileId] = struct{}{}
	}

	// attachments.file_id
	type attachRow struct {
		FileId int `orm:"file_id"`
	}
	var attachments []attachRow
	g.DB().Ctx(ctx).Model("attachments").
		Fields("file_id").
		Scan(&attachments)
	for _, a := range attachments {
		referenced[a.FileId] = struct{}{}
	}

	return referenced
}

// getFileDiskPath 根据 hash 和文件名计算磁盘路径
func getFileDiskPath(ctx context.Context, hash []byte, name string) string {
	uploadPath := g.Cfg().MustGet(ctx, "upload.path").String()
	if uploadPath == "" {
		return ""
	}
	hashStr := files.HashToString(hash)
	ext := ""
	if dotIdx := gstr.StrEx(name, "."); dotIdx != "" {
		ext = "." + dotIdx
	}
	fileName := hashStr + ext
	return filepath.Join(uploadPath, fileName[0:2], fileName[2:4], fileName)
}

// cleanEmptyDirs 清理空的父目录
func cleanEmptyDirs(filePath string) {
	parent := filepath.Dir(filePath)
	entries, err := os.ReadDir(parent)
	if err == nil && len(entries) == 0 {
		os.Remove(parent)
		grandparent := filepath.Dir(parent)
		entries2, err2 := os.ReadDir(grandparent)
		if err2 == nil && len(entries2) == 0 {
			os.Remove(grandparent)
		}
	}
}
