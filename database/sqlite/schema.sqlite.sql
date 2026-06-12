-- ==============================================================
-- SuAsk SQLite Schema (v3, post migration 003)
--
-- 本文件是"从零新建 DB"用的；已有 DB 请走 migrations/ 里按编号升级。
--
-- 相比 v2 的变化：
--   - 去掉了 "问大家" 模块（questions.dst_user_id 现在 NOT NULL）
--   - 为删除业务铺路：questions/answers/files 加 deleted_at 软删字段；
--     全部外键按"依附关系 CASCADE、主数据 SET NULL"重设
--   - 为公告模块铺路：answers/attachments 加 announcement_id（公告评论/公告附件）
--   - 补充部分索引（仅索引未软删行），以及若干组合索引
-- ==============================================================

PRAGMA foreign_keys = OFF;
PRAGMA journal_mode = WAL;

-- ------------------------- users -------------------------
CREATE TABLE users (
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
CREATE INDEX idx_users_avatar_file ON users(avatar_file_id);

-- ------------------------- guest_users -------------------------
-- 记录 guest 用户的过期时间，过期后由定时任务硬删除
CREATE TABLE guest_users (
  id          INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  expires_at  DATETIME NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_guest_users_expires ON guest_users(expires_at);

-- ------------------------- files -------------------------
CREATE TABLE files (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT    NOT NULL,
  hash        BLOB    NOT NULL,
  uploader_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
  created_at  DATETIME,
  deleted_at  DATETIME
);
CREATE INDEX idx_files_uploader ON files(uploader_id);
CREATE INDEX idx_files_alive ON files(id) WHERE deleted_at IS NULL;

-- ------------------------- teachers -------------------------
CREATE TABLE teachers (
  id           INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  responses    INTEGER DEFAULT 0,
  name         TEXT,
  avatar_url   TEXT,
  introduction TEXT,
  email        TEXT,
  perm         TEXT CHECK (perm IN ('public','protected','private'))
);

-- ------------------------- settings -------------------------
CREATE TABLE settings (
  id                 INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  theme_id           INTEGER,
  question_box_perm  TEXT CHECK (question_box_perm IN ('public','protected','private')),
  notify_switch      INTEGER NOT NULL DEFAULT 1 CHECK (notify_switch IN (0,1)),
  notify_email       TEXT,
  notify_merge_cnt   INTEGER DEFAULT 1,
  notify_max_delay   INTEGER DEFAULT 0
);

-- ------------------------- themes / config -------------------------
CREATE TABLE themes (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  background_path TEXT NOT NULL
);

CREATE TABLE config (
  id                  INTEGER PRIMARY KEY CHECK (id = 0),
  default_avatar_path TEXT NOT NULL,
  default_theme_id    INTEGER NOT NULL
);

-- ------------------------- questions -------------------------
-- dst_user_id 现在必填：所有问题都必须指定目标老师。
-- 用户注销（软删）时 src/dst 会变 NULL（ON DELETE SET NULL），
-- 问题本身保留，只是"发起人"或"接收者"变成匿名/无主。
CREATE TABLE questions (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  src_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
  dst_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
  title       TEXT    NOT NULL,
  contents    TEXT    NOT NULL,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  views       INTEGER NOT NULL DEFAULT 0,
  reply_cnt   INTEGER NOT NULL DEFAULT 0,
  deleted_at  DATETIME
);
CREATE INDEX idx_questions_src         ON questions(src_user_id);
CREATE INDEX idx_questions_dst         ON questions(dst_user_id);
CREATE INDEX idx_questions_created_at  ON questions(created_at);
CREATE INDEX idx_questions_dst_reply   ON questions(dst_user_id, reply_cnt);
CREATE INDEX idx_questions_alive_dst_created
  ON questions(dst_user_id, created_at DESC)
  WHERE deleted_at IS NULL;

-- ------------------------- answers -------------------------
-- 一条回答要么挂在 question 下（普通问题的回答），
-- 要么挂在 announcement 下（公告评论）。XOR 约束强制二选一。
-- announcement_id 的外键约束放在 migration 004（announcements 表建好之后）
-- 一起补上；在那之前先留 INTEGER 占位。
CREATE TABLE answers (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id         INTEGER REFERENCES users(id)     ON DELETE SET NULL,
  question_id     INTEGER REFERENCES questions(id) ON DELETE CASCADE,
  announcement_id INTEGER,  -- TODO 在 migration 004 补 REFERENCES announcements(id) ON DELETE CASCADE
  in_reply_to     INTEGER REFERENCES answers(id)   ON DELETE SET NULL,
  contents        TEXT    NOT NULL,
  created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  upvotes         INTEGER NOT NULL DEFAULT 0,
  deleted_at      DATETIME,
  CHECK ((question_id IS NULL) <> (announcement_id IS NULL))
);
CREATE INDEX idx_answers_question        ON answers(question_id);
CREATE INDEX idx_answers_user            ON answers(user_id);
CREATE INDEX idx_answers_in_reply_to     ON answers(in_reply_to);
CREATE INDEX idx_answers_announcement    ON answers(announcement_id);
CREATE INDEX idx_answers_alive_question
  ON answers(question_id, created_at)
  WHERE deleted_at IS NULL AND question_id IS NOT NULL;
CREATE INDEX idx_answers_alive_announcement
  ON answers(announcement_id, created_at)
  WHERE deleted_at IS NULL AND announcement_id IS NOT NULL;

-- ------------------------- attachments -------------------------
-- 三选一：question / answer / announcement。
CREATE TABLE attachments (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  question_id     INTEGER REFERENCES questions(id) ON DELETE CASCADE,
  answer_id       INTEGER REFERENCES answers(id)   ON DELETE CASCADE,
  announcement_id INTEGER,  -- TODO migration 004
  type            TEXT    NOT NULL CHECK (type IN ('picture')),
  file_id         INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  CHECK (
    (CASE WHEN question_id     IS NULL THEN 0 ELSE 1 END) +
    (CASE WHEN answer_id       IS NULL THEN 0 ELSE 1 END) +
    (CASE WHEN announcement_id IS NULL THEN 0 ELSE 1 END) = 1
  )
);
CREATE INDEX idx_attach_question     ON attachments(question_id);
CREATE INDEX idx_attach_answer       ON attachments(answer_id);
CREATE INDEX idx_attach_announcement ON attachments(announcement_id);
CREATE INDEX idx_attach_file         ON attachments(file_id);

-- ------------------------- favorites -------------------------
-- package='default' 学生收藏；package='top' 教师置顶。
CREATE TABLE favorites (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  package     TEXT    NOT NULL DEFAULT 'default'
                       CHECK (package IN ('default','top')),
  UNIQUE (user_id, question_id, package)
);
CREATE INDEX idx_fav_question     ON favorites(question_id);
CREATE INDEX idx_fav_user_package ON favorites(user_id, package);

-- ------------------------- upvotes -------------------------
CREATE TABLE upvotes (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER          REFERENCES questions(id) ON DELETE CASCADE,
  answer_id   INTEGER          REFERENCES answers(id)   ON DELETE CASCADE,
  CHECK ((question_id IS NULL) <> (answer_id IS NULL))
);
CREATE INDEX idx_upvotes_user        ON upvotes(user_id);
CREATE INDEX idx_upvotes_answer      ON upvotes(answer_id);
CREATE INDEX idx_upvotes_user_answer ON upvotes(user_id, answer_id);

-- ------------------------- notifications -------------------------
CREATE TABLE notifications (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  reply_to_id INTEGER          REFERENCES answers(id)   ON DELETE CASCADE,
  answer_id   INTEGER          REFERENCES answers(id)   ON DELETE CASCADE,
  type        TEXT    NOT NULL CHECK (type IN ('new_question','new_reply','new_answer')),
  is_read     INTEGER NOT NULL DEFAULT 0 CHECK (is_read IN (0,1)),
  created_at  DATETIME,
  deleted_at  DATETIME
);
CREATE INDEX idx_notif_user_read_type
  ON notifications(user_id, is_read, type, created_at);
CREATE INDEX idx_notif_question ON notifications(question_id);

-- ------------------------- user_relation -------------------------
-- 记录每个问题的前 N 个回复者，用于列表页头像拼接
CREATE TABLE user_relation (
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  PRIMARY KEY (question_id, user_id)
);

-- ------------------------- viewed -------------------------
CREATE TABLE viewed (
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, question_id)
);

PRAGMA foreign_keys = ON;
