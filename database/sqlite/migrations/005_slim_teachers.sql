-- ============================================================================
-- Migration 005: Slim teachers table
--
-- teachers 表原来冗余了 name/email/introduction/avatar_url，和 users 表重复。
-- 精简为只保留 id, perm, responses。
-- 老师的 name/email/introduction/avatar 统一从 users 表读取。
--
-- 注意：teachers.avatar_url 是外部链接（不是 file_id），保留到 users.avatar_file_id
-- 的迁移需要人工处理（把外部 URL 下载为文件再入库），这里暂不自动迁移，
-- 只是把 avatar_url 列删掉。前端改为优先读 users.avatar_file_id，
-- 如果为空则 fallback 到 teachers.avatar_url（通过新接口返回）。
--
-- 实际做法：保留 avatar_url 列（兼容期），但标记为 deprecated。
-- ============================================================================

PRAGMA foreign_keys = OFF;

BEGIN;

-- 重建 teachers 表：只保留核心字段 + avatar_url（兼容期保留）
ALTER TABLE teachers RENAME TO teachers_old_005;

CREATE TABLE teachers (
  id           INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  perm         TEXT CHECK (perm IN ('public','protected','private')),
  responses    INTEGER DEFAULT 0,
  avatar_url   TEXT  -- deprecated: 兼容期保留，前端迁移完后可删
);

INSERT INTO teachers (id, perm, responses, avatar_url)
  SELECT id, perm, responses, avatar_url FROM teachers_old_005;

DROP TABLE teachers_old_005;

COMMIT;

PRAGMA foreign_keys = ON;
