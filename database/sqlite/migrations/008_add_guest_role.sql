-- Migration 008: 添加 guest 角色支持
-- 
-- 变更说明：
--   1. 重建 users 表：放宽 email/salt/password_hash 的 NOT NULL 约束，role CHECK 加入 'guest'
--   2. 新建 guest_users 副表：记录 guest 用户的过期时间
--   3. 添加 idx_guest_users_expires 索引
--
-- SQLite 不支持 ALTER COLUMN / ALTER CHECK，必须重建表

PRAGMA foreign_keys = OFF;

-- 1. 重建 users 表
CREATE TABLE users_new (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  name           TEXT    NOT NULL UNIQUE,
  email          TEXT    UNIQUE,
  salt           TEXT,
  password_hash  TEXT,
  role           TEXT    NOT NULL CHECK (role IN ('admin','teacher','student','guest')),
  nickname       TEXT    NOT NULL,
  introduction   TEXT    NOT NULL DEFAULT '',
  avatar_file_id INTEGER,
  created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     DATETIME,
  deleted_at     DATETIME
);

INSERT INTO users_new (id, name, email, salt, password_hash, role, nickname, introduction, avatar_file_id, created_at, updated_at, deleted_at)
  SELECT id, name, email, salt, password_hash, role, nickname, introduction, avatar_file_id, created_at, updated_at, deleted_at FROM users;

DROP TABLE users;
ALTER TABLE users_new RENAME TO users;

-- 重建索引
CREATE INDEX idx_users_avatar_file ON users(avatar_file_id);

-- 2. 创建 guest_users 副表
CREATE TABLE guest_users (
  id          INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  expires_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. 添加过期时间索引
CREATE INDEX idx_guest_users_expires ON guest_users(expires_at);

PRAGMA foreign_keys = ON;
