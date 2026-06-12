-- ============================================================================
-- Migration 003: Drop "问大家" + prepare for delete/announcement modules
--
-- 本脚本做以下事：
--   1) 硬删所有 dst_user_id IS NULL 的 questions 及其所有级联（手动按依赖顺序删，
--      因为当前外键还是 RESTRICT，不能一条 DELETE 解决）
--   2) 重建 4 张核心表：questions / answers / attachments / files
--      - 外键从 RESTRICT 改为按"依附性 vs 主数据"分类：CASCADE / SET NULL / RESTRICT
--      - questions / answers / files 加 deleted_at 软删字段
--      - answers 加 announcement_id（公告评论复用）
--      - attachments 加 announcement_id + CHECK 改为三选一（question/answer/announcement）
--      - questions.dst_user_id 改为 NOT NULL（"问大家"已清除，所有问题必须有目标老师）
--   3) 其余 8 张外键从 users/questions/answers 出发的表，重建以把 RESTRICT
--      改成 CASCADE（收藏/点赞/通知/user_relation/viewed 这些都是依附关系，
--      父删了就没意义）
--   4) 重建索引，新增部分索引（只索引未删行）
--
-- 前提：
--   - 本脚本用于"同步生产 DB 之前的测试库"；跑前先备份
--   - 执行过程用事务包住，失败自动回滚
-- ============================================================================

PRAGMA foreign_keys = OFF;

BEGIN;

-- --------------------------------------------------------------------------
-- Step 1: 硬删"问大家"数据及其全部级联
--    当前 FK 是 RESTRICT，必须按依赖从叶子删到根，否则会报 FOREIGN KEY constraint failed
-- --------------------------------------------------------------------------

-- 1.1 收集待删的 question_id 和 answer_id 到临时表，方便后续重复用
CREATE TEMP TABLE _del_q AS
  SELECT id FROM questions WHERE dst_user_id IS NULL;
CREATE TEMP TABLE _del_a AS
  SELECT id FROM answers WHERE question_id IN (SELECT id FROM _del_q);

-- 1.2 attachments：answer 的附件先删，再删 question 的附件
DELETE FROM attachments WHERE answer_id   IN (SELECT id FROM _del_a);
DELETE FROM attachments WHERE question_id IN (SELECT id FROM _del_q);

-- 1.3 upvotes：涉及被删 answer 或 question 的
DELETE FROM upvotes WHERE answer_id   IN (SELECT id FROM _del_a);
DELETE FROM upvotes WHERE question_id IN (SELECT id FROM _del_q);

-- 1.4 notifications：涉及被删 question 或 answer 的
DELETE FROM notifications WHERE answer_id   IN (SELECT id FROM _del_a);
DELETE FROM notifications WHERE reply_to_id IN (SELECT id FROM _del_a);
DELETE FROM notifications WHERE question_id IN (SELECT id FROM _del_q);

-- 1.5 favorites / user_relation / viewed
DELETE FROM favorites     WHERE question_id IN (SELECT id FROM _del_q);
DELETE FROM user_relation WHERE question_id IN (SELECT id FROM _del_q);
DELETE FROM viewed        WHERE question_id IN (SELECT id FROM _del_q);

-- 1.6 answers in_reply_to 链内部可能互相引用：先把 in_reply_to 设 NULL 解套
UPDATE answers SET in_reply_to = NULL
  WHERE in_reply_to IN (SELECT id FROM _del_a);

-- 1.7 最后删 answers 再删 questions
DELETE FROM answers   WHERE id IN (SELECT id FROM _del_a);
DELETE FROM questions WHERE id IN (SELECT id FROM _del_q);

DROP TABLE _del_a;
DROP TABLE _del_q;

-- --------------------------------------------------------------------------
-- Step 2: 重建 questions（加 deleted_at；dst_user_id 改为 NOT NULL；FK 重设）
-- --------------------------------------------------------------------------
CREATE TABLE questions_new (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  src_user_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
  dst_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE SET NULL,
  title       TEXT    NOT NULL,
  contents    TEXT    NOT NULL,
  is_private  INTEGER NOT NULL DEFAULT 0 CHECK (is_private IN (0,1)),
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  views       INTEGER NOT NULL DEFAULT 0,
  reply_cnt   INTEGER NOT NULL DEFAULT 0,
  deleted_at  DATETIME
);
INSERT INTO questions_new
  (id, src_user_id, dst_user_id, title, contents, is_private, created_at, views, reply_cnt)
  SELECT id, src_user_id, dst_user_id, title, contents, is_private, created_at, views, reply_cnt
    FROM questions;
DROP TABLE questions;
ALTER TABLE questions_new RENAME TO questions;

-- --------------------------------------------------------------------------
-- Step 3: 重建 answers（加 deleted_at + announcement_id；FK 重设）
--    announcement_id 现在还没 announcements 表，先用 INTEGER 占位，
--    announcements 建好后用下一个 migration (004) 加外键约束。
--
--    in_reply_to 是自引用外键；SQLite 会把 REFERENCES 的表名按字面存，
--    所以这里直接用最终名字 answers：先把旧表改名、再 CREATE 新 answers
--    + 复制数据、再删掉改名后的旧表。
-- --------------------------------------------------------------------------
ALTER TABLE answers RENAME TO answers_old;
CREATE TABLE answers (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id         INTEGER REFERENCES users(id)     ON DELETE SET NULL,
  question_id     INTEGER          REFERENCES questions(id) ON DELETE CASCADE,
  announcement_id INTEGER,  -- 将在 migration 004 里加 REFERENCES announcements(id) ON DELETE CASCADE
  in_reply_to     INTEGER          REFERENCES answers(id) ON DELETE SET NULL,
  contents        TEXT    NOT NULL,
  created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  upvotes         INTEGER NOT NULL DEFAULT 0,
  deleted_at      DATETIME,
  -- 一条回答要么挂在 question 下，要么挂在 announcement 下，XOR
  CHECK ((question_id IS NULL) <> (announcement_id IS NULL))
);
INSERT INTO answers
  (id, user_id, question_id, announcement_id, in_reply_to, contents, created_at, upvotes)
  SELECT id, user_id, question_id, NULL, in_reply_to, contents, created_at, upvotes
    FROM answers_old;
DROP TABLE answers_old;

-- --------------------------------------------------------------------------
-- Step 4: 重建 attachments（加 announcement_id + CHECK 三选一；FK 重设）
-- --------------------------------------------------------------------------
CREATE TABLE attachments_new (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  question_id     INTEGER REFERENCES questions(id)     ON DELETE CASCADE,
  answer_id       INTEGER REFERENCES answers(id)       ON DELETE CASCADE,
  announcement_id INTEGER,  -- 同 answers，等 004 migration 再加外键
  type            TEXT    NOT NULL CHECK (type IN ('picture')),
  file_id         INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  -- question_id / answer_id / announcement_id 必须恰好一个非 NULL
  CHECK (
    (CASE WHEN question_id     IS NULL THEN 0 ELSE 1 END) +
    (CASE WHEN answer_id       IS NULL THEN 0 ELSE 1 END) +
    (CASE WHEN announcement_id IS NULL THEN 0 ELSE 1 END) = 1
  )
);
INSERT INTO attachments_new
  (id, question_id, answer_id, announcement_id, type, file_id)
  SELECT id, question_id, answer_id, NULL, type, file_id
    FROM attachments;
DROP TABLE attachments;
ALTER TABLE attachments_new RENAME TO attachments;

-- --------------------------------------------------------------------------
-- Step 5: 重建 files（加 deleted_at；uploader_id 改 SET NULL）
-- --------------------------------------------------------------------------
CREATE TABLE files_new (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        TEXT    NOT NULL,
  hash        BLOB    NOT NULL,
  uploader_id INTEGER REFERENCES users(id) ON DELETE SET NULL,
  created_at  DATETIME,
  deleted_at  DATETIME
);
INSERT INTO files_new (id, name, hash, uploader_id, created_at)
  SELECT id, name, hash, uploader_id, created_at FROM files;
DROP TABLE files;
ALTER TABLE files_new RENAME TO files;

-- --------------------------------------------------------------------------
-- Step 6: 重建 favorites / upvotes / notifications / user_relation / viewed
--    只为了把 FK 从 RESTRICT 改成 CASCADE
-- --------------------------------------------------------------------------
CREATE TABLE favorites_new (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  created_at  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  package     TEXT NOT NULL DEFAULT 'default' CHECK (package IN ('default','top')),
  UNIQUE (user_id, question_id, package)
);
INSERT INTO favorites_new SELECT * FROM favorites;
DROP TABLE favorites;
ALTER TABLE favorites_new RENAME TO favorites;

CREATE TABLE upvotes_new (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER          REFERENCES questions(id) ON DELETE CASCADE,
  answer_id   INTEGER          REFERENCES answers(id)   ON DELETE CASCADE,
  CHECK ((question_id IS NULL) <> (answer_id IS NULL))
);
INSERT INTO upvotes_new SELECT * FROM upvotes;
DROP TABLE upvotes;
ALTER TABLE upvotes_new RENAME TO upvotes;

CREATE TABLE notifications_new (
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
INSERT INTO notifications_new
  (id, user_id, question_id, reply_to_id, answer_id, type, is_read, created_at, deleted_at)
  SELECT id, user_id, question_id, reply_to_id, answer_id, type, is_read, created_at, deleted_at
    FROM notifications;
DROP TABLE notifications;
ALTER TABLE notifications_new RENAME TO notifications;

CREATE TABLE user_relation_new (
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  PRIMARY KEY (question_id, user_id)
);
INSERT INTO user_relation_new SELECT * FROM user_relation;
DROP TABLE user_relation;
ALTER TABLE user_relation_new RENAME TO user_relation;

CREATE TABLE viewed_new (
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
  PRIMARY KEY (user_id, question_id)
);
INSERT INTO viewed_new SELECT * FROM viewed;
DROP TABLE viewed;
ALTER TABLE viewed_new RENAME TO viewed;

-- --------------------------------------------------------------------------
-- Step 7: 重建索引（原 schema 里有的全部重来 + 补充部分索引）
-- 先把 users 上的已有索引删掉，避免"already exists"。
-- questions/answers/attachments/files/favorites/upvotes/notifications 因为
-- 表刚被 DROP 重建，索引自动就没了，无需 DROP。
-- --------------------------------------------------------------------------

DROP INDEX IF EXISTS idx_users_avatar_file;

-- questions
CREATE INDEX idx_questions_src           ON questions(src_user_id);
CREATE INDEX idx_questions_dst           ON questions(dst_user_id);
CREATE INDEX idx_questions_created_at    ON questions(created_at);
CREATE INDEX idx_questions_dst_reply     ON questions(dst_user_id, reply_cnt);
-- 部分索引：只覆盖未软删的问题，列表查询走这里就完全不碰已删行
CREATE INDEX idx_questions_alive_dst_created
  ON questions(dst_user_id, created_at DESC)
  WHERE deleted_at IS NULL;

-- answers
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

-- attachments
CREATE INDEX idx_attach_question     ON attachments(question_id);
CREATE INDEX idx_attach_answer       ON attachments(answer_id);
CREATE INDEX idx_attach_announcement ON attachments(announcement_id);
CREATE INDEX idx_attach_file         ON attachments(file_id);

-- files
CREATE INDEX idx_files_uploader ON files(uploader_id);
CREATE INDEX idx_files_alive
  ON files(id) WHERE deleted_at IS NULL;

-- users（avatar_file_id 的辅助索引）
CREATE INDEX idx_users_avatar_file ON users(avatar_file_id);

-- favorites
CREATE INDEX idx_fav_question ON favorites(question_id);
-- 新增：大部分 favorites 查询是 (user_id, package) 组合
CREATE INDEX idx_fav_user_package ON favorites(user_id, package);

-- upvotes
CREATE INDEX idx_upvotes_user    ON upvotes(user_id);
CREATE INDEX idx_upvotes_answer  ON upvotes(answer_id);
-- 新增：详情页常用"我是否点过这个 answer 的赞"的组合查询
CREATE INDEX idx_upvotes_user_answer ON upvotes(user_id, answer_id);

-- notifications
-- 未读计数按 (user_id, is_read, type) 过滤，is_read 基数比 type 大，放中间
CREATE INDEX idx_notif_user_read_type
  ON notifications(user_id, is_read, type, created_at);
CREATE INDEX idx_notif_question ON notifications(question_id);

COMMIT;

PRAGMA foreign_keys = ON;
-- 跑完必须再单独执行 `PRAGMA foreign_key_check;`，看没有任何输出才算安全
