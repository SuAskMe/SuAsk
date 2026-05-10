-- ==============================================================
-- SuAsk SQLite Schema
-- 由 MySQL 8.0 schema 迁移而来，去除 MySQL 独有语法与类型
-- 迁移指南：见同目录下 README.md 与 migrate.py
-- ==============================================================

PRAGMA foreign_keys = OFF;   -- 建表阶段关闭外键（顺序无关），灌数据前再打开
PRAGMA journal_mode = WAL;

-- ------------------------- users -------------------------
CREATE TABLE users (
  id             INTEGER PRIMARY KEY AUTOINCREMENT,
  name           TEXT    NOT NULL UNIQUE,          -- 用户名唯一
  email          TEXT    NOT NULL UNIQUE,          -- 邮箱唯一
  salt           TEXT    NOT NULL,                 -- 加密盐
  password_hash  TEXT    NOT NULL,                 -- 密码哈希
  role           TEXT    NOT NULL CHECK (role IN ('admin','teacher','student')),
  nickname       TEXT    NOT NULL,
  introduction   TEXT    NOT NULL DEFAULT '',
  avatar_file_id INTEGER,                          -- FK 在 files 建表后通过外键约束声明
  created_at     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at     DATETIME,
  deleted_at     DATETIME
);

-- ------------------------- files -------------------------
CREATE TABLE files (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT    NOT NULL,
  hash        BLOB    NOT NULL,                    -- 32B BLAKE2b
  uploader_id INTEGER REFERENCES users(id) ON DELETE RESTRICT,
  created_at  DATETIME
);
CREATE INDEX idx_files_uploader ON files(uploader_id);

-- users.avatar_file_id 的外键用一个 trigger 代替（SQLite 不支持后加 FK，现在加索引足够）
CREATE INDEX idx_users_avatar_file ON users(avatar_file_id);

-- ------------------------- teachers -------------------------
CREATE TABLE teachers (
  id           INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE RESTRICT,
  responses    INTEGER DEFAULT 0,
  name         TEXT,
  avatar_url   TEXT,
  introduction TEXT,
  email        TEXT,
  perm         TEXT CHECK (perm IN ('public','protected','private'))
);

-- ------------------------- settings -------------------------
CREATE TABLE settings (
  id                 INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE RESTRICT,
  theme_id           INTEGER,
  question_box_perm  TEXT CHECK (question_box_perm IN ('public','protected','private')),
  notify_switch      INTEGER NOT NULL DEFAULT 1 CHECK (notify_switch IN (0,1)),
  notify_email       TEXT,
  notify_merge_cnt   INTEGER DEFAULT 1,
  notify_max_delay   INTEGER DEFAULT 0
);

-- ------------------------- themes -------------------------
CREATE TABLE themes (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  background_path TEXT NOT NULL
);

-- ------------------------- config -------------------------
-- 单行配置：强制 id = 0
CREATE TABLE config (
  id                  INTEGER PRIMARY KEY CHECK (id = 0),
  default_avatar_path TEXT NOT NULL,
  default_theme_id    INTEGER NOT NULL
);

-- ------------------------- questions -------------------------
CREATE TABLE questions (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  src_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  dst_user_id INTEGER          REFERENCES users(id) ON DELETE RESTRICT,
  title       TEXT    NOT NULL,
  contents    TEXT    NOT NULL,
  is_private  INTEGER NOT NULL DEFAULT 0 CHECK (is_private IN (0,1)),
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  views       INTEGER NOT NULL DEFAULT 0,
  reply_cnt   INTEGER NOT NULL DEFAULT 0,
  -- 私密提问必须指定教师
  CHECK (dst_user_id IS NOT NULL OR is_private = 0)
);
CREATE INDEX idx_questions_src        ON questions(src_user_id);
CREATE INDEX idx_questions_dst        ON questions(dst_user_id);
CREATE INDEX idx_questions_created_at ON questions(created_at);
-- 用于 LIKE 搜索的辅助索引（SQLite 对 '%foo%' 这种前后通配无法走索引，这里只为了 reply_cnt 过滤）
CREATE INDEX idx_questions_dst_reply  ON questions(dst_user_id, reply_cnt);

-- ------------------------- answers -------------------------
CREATE TABLE answers (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE RESTRICT,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
  in_reply_to INTEGER          REFERENCES answers(id)   ON DELETE RESTRICT,
  contents    TEXT    NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  upvotes     INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_answers_question    ON answers(question_id);
CREATE INDEX idx_answers_user        ON answers(user_id);
CREATE INDEX idx_answers_in_reply_to ON answers(in_reply_to);

-- ------------------------- attachments -------------------------
CREATE TABLE attachments (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  question_id INTEGER REFERENCES questions(id) ON DELETE RESTRICT,
  answer_id   INTEGER REFERENCES answers(id)   ON DELETE RESTRICT,
  type        TEXT    NOT NULL CHECK (type IN ('picture')),
  file_id     INTEGER NOT NULL REFERENCES files(id) ON DELETE RESTRICT,
  -- question_id XOR answer_id 必须恰好有一个
  CHECK ((question_id IS NULL) <> (answer_id IS NULL))
);
CREATE INDEX idx_attach_question ON attachments(question_id);
CREATE INDEX idx_attach_answer   ON attachments(answer_id);
CREATE INDEX idx_attach_file     ON attachments(file_id);

-- ------------------------- favorites -------------------------
-- 同一张表承载两种语义：
--   package = 'default' → 学生收藏
--   package = 'top'     → 教师置顶（原 MySQL 里用 '置顶'，这里英文化；迁移脚本会自动转换）
CREATE TABLE favorites (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE RESTRICT,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  package     TEXT    NOT NULL DEFAULT 'default'
                       CHECK (package IN ('default','top')),
  UNIQUE (user_id, question_id, package)
);
CREATE INDEX idx_fav_question ON favorites(question_id);

-- ------------------------- upvotes -------------------------
CREATE TABLE upvotes (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE RESTRICT,
  question_id INTEGER          REFERENCES questions(id) ON DELETE RESTRICT,
  answer_id   INTEGER          REFERENCES answers(id)   ON DELETE RESTRICT,
  CHECK ((question_id IS NULL) <> (answer_id IS NULL))
);
CREATE INDEX idx_upvotes_user   ON upvotes(user_id);
CREATE INDEX idx_upvotes_answer ON upvotes(answer_id);

-- ------------------------- notifications -------------------------
CREATE TABLE notifications (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE RESTRICT,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
  reply_to_id INTEGER REFERENCES answers(id) ON DELETE RESTRICT,
  answer_id   INTEGER REFERENCES answers(id) ON DELETE RESTRICT,
  type        TEXT    NOT NULL CHECK (type IN ('new_question','new_reply','new_answer')),
  is_read     INTEGER NOT NULL DEFAULT 0 CHECK (is_read IN (0,1)),
  created_at  DATETIME,
  deleted_at  DATETIME
);
CREATE INDEX idx_notif_user_type    ON notifications(user_id, type, is_read, created_at);
CREATE INDEX idx_notif_question     ON notifications(question_id);

-- ------------------------- user_relation -------------------------
-- 记录每个问题的前 N 个回复者，用于列表页头像拼接
CREATE TABLE user_relation (
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE RESTRICT,
  PRIMARY KEY (question_id, user_id)
);

-- ------------------------- viewed -------------------------
-- 目前未被业务代码引用，保留以备后续"去重浏览量"使用
CREATE TABLE viewed (
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE RESTRICT,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, question_id)
);

-- 全部建完再开外键
PRAGMA foreign_keys = ON;
