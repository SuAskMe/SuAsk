-- Migration 007: 删除 teachers.avatar_url 列
-- 前置条件: 先运行 007_migrate_teacher_avatars.py 完成头像迁移
--
-- 教师头像统一使用 users.avatar_file_id，不再需要 teachers.avatar_url

-- SQLite 3.35.0+ 支持 ALTER TABLE DROP COLUMN
ALTER TABLE teachers DROP COLUMN avatar_url;
