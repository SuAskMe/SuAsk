-- ============================================================================
-- Migration 004: Create announcements table + add FK on answers/attachments
--
-- 前提：003 已跑完（answers/attachments 已有 announcement_id 列但无外键约束）
-- ============================================================================

PRAGMA foreign_keys = OFF;

BEGIN;

-- 1) 创建 announcements 表
CREATE TABLE IF NOT EXISTS announcements (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  author_id    INTEGER REFERENCES users(id) ON DELETE SET NULL,
  title        TEXT    NOT NULL,
  contents     TEXT    NOT NULL,                       -- Markdown
  is_pinned    INTEGER NOT NULL DEFAULT 0 CHECK (is_pinned IN (0,1)),
  published_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expires_at   DATETIME,                               -- NULL = 永不过期
  updated_at   DATETIME,
  deleted_at   DATETIME
);

CREATE INDEX IF NOT EXISTS idx_announcements_live
  ON announcements(is_pinned DESC, published_at DESC)
  WHERE deleted_at IS NULL;

-- 2) 重建 answers 以加上 announcement_id 的外键约束
--    （003 里只是 INTEGER 占位，没有 REFERENCES）
ALTER TABLE answers RENAME TO answers_old_004;

CREATE TABLE answers (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id         INTEGER REFERENCES users(id)          ON DELETE SET NULL,
  question_id     INTEGER REFERENCES questions(id)      ON DELETE CASCADE,
  announcement_id INTEGER REFERENCES announcements(id)  ON DELETE CASCADE,
  in_reply_to     INTEGER REFERENCES answers(id)        ON DELETE SET NULL,
  contents        TEXT    NOT NULL,
  created_at      DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  upvotes         INTEGER NOT NULL DEFAULT 0,
  deleted_at      DATETIME,
  CHECK ((question_id IS NULL) <> (announcement_id IS NULL))
);

INSERT INTO answers
  (id, user_id, question_id, announcement_id, in_reply_to, contents, created_at, upvotes, deleted_at)
  SELECT id, user_id, question_id, announcement_id, in_reply_to, contents, created_at, upvotes, deleted_at
    FROM answers_old_004;

DROP TABLE answers_old_004;

-- 3) 重建 attachments 以加上 announcement_id 的外键约束
ALTER TABLE attachments RENAME TO attachments_old_004;

CREATE TABLE attachments (
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  question_id     INTEGER REFERENCES questions(id)      ON DELETE CASCADE,
  answer_id       INTEGER REFERENCES answers(id)        ON DELETE CASCADE,
  announcement_id INTEGER REFERENCES announcements(id)  ON DELETE CASCADE,
  type            TEXT    NOT NULL CHECK (type IN ('picture')),
  file_id         INTEGER NOT NULL REFERENCES files(id) ON DELETE CASCADE,
  CHECK (
    (CASE WHEN question_id     IS NULL THEN 0 ELSE 1 END) +
    (CASE WHEN answer_id       IS NULL THEN 0 ELSE 1 END) +
    (CASE WHEN announcement_id IS NULL THEN 0 ELSE 1 END) = 1
  )
);

INSERT INTO attachments
  (id, question_id, answer_id, announcement_id, type, file_id)
  SELECT id, question_id, answer_id, announcement_id, type, file_id
    FROM attachments_old_004;

DROP TABLE attachments_old_004;

-- 4) Rebuild notifications to fix FK references (after answers table rebuild)
ALTER TABLE notifications RENAME TO notifications_old_004;

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

INSERT INTO notifications
  (id, user_id, question_id, reply_to_id, answer_id, type, is_read, created_at, deleted_at)
  SELECT id, user_id, question_id, reply_to_id, answer_id, type, is_read, created_at, deleted_at
    FROM notifications_old_004;

DROP TABLE notifications_old_004;

CREATE INDEX idx_notif_user_read_type
  ON notifications(user_id, is_read, type, created_at);
CREATE INDEX idx_notif_question ON notifications(question_id);

-- 5) Also rebuild upvotes (references answers)
ALTER TABLE upvotes RENAME TO upvotes_old_004;

CREATE TABLE upvotes (
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id     INTEGER NOT NULL REFERENCES users(id)     ON DELETE CASCADE,
  question_id INTEGER          REFERENCES questions(id) ON DELETE CASCADE,
  answer_id   INTEGER          REFERENCES answers(id)   ON DELETE CASCADE,
  CHECK ((question_id IS NULL) <> (answer_id IS NULL))
);

INSERT INTO upvotes SELECT * FROM upvotes_old_004;
DROP TABLE upvotes_old_004;

CREATE INDEX idx_upvotes_user        ON upvotes(user_id);
CREATE INDEX idx_upvotes_answer      ON upvotes(answer_id);
CREATE INDEX idx_upvotes_user_answer ON upvotes(user_id, answer_id);

-- 6) Rebuild favorites (references questions which didn't change, but let's be safe)
--    Actually favorites FK to questions is fine. Skip.

-- 7) Reindex answers + attachments
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

CREATE INDEX idx_attach_question     ON attachments(question_id);
CREATE INDEX idx_attach_answer       ON attachments(answer_id);
CREATE INDEX idx_attach_announcement ON attachments(announcement_id);
CREATE INDEX idx_attach_file         ON attachments(file_id);

COMMIT;

PRAGMA foreign_keys = ON;
