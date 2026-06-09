"""
给回归测试准备一个确定的用户账号，直接写 SQLite（绕过邮件验证码）。

用户名/密码默认是 smoke_student / Smoke-Pass-123! ，可用 CLI 参数覆盖。

密码哈希策略：
  - 新账号直接写 bcrypt，和 utility.HashPassword 保持一致
  - 如果机器上没装 `bcrypt` python 包，回退到老的 md5(md5(pw)+md5(salt))，后端登录
    时会透明升级到 bcrypt，所以 smoke 仍然能过

一旦后端改了算法（比如上 argon2id），这里也要同步；这就是我们要的"提醒"。
"""

from __future__ import annotations

import hashlib
import sqlite3
from pathlib import Path

try:
    import bcrypt  # pip install bcrypt
    _HAS_BCRYPT = True
except ImportError:
    _HAS_BCRYPT = False


def md5_hex(s: str) -> str:
    return hashlib.md5(s.encode("utf-8")).hexdigest()


def legacy_encrypt_password(password: str, salt: str) -> str:
    """与 utility.EncryptPassword（老算法）保持一致，仅用于无 bcrypt 环境兜底。"""
    return md5_hex(md5_hex(password) + md5_hex(salt))


def hash_password(password: str) -> tuple[str, str]:
    """返回 (salt, hash)；bcrypt 走 salt=空。"""
    if _HAS_BCRYPT:
        h = bcrypt.hashpw(password.encode("utf-8"), bcrypt.gensalt(rounds=10))
        return "", h.decode("utf-8")
    # 兜底：用老算法写进去，登录后后端会透明升级
    salt = "smokesalt1"
    return salt, legacy_encrypt_password(password, salt)


SMOKE_QUESTION_SEEDS = [
    {
        "title": "[smoke] teacher inbox unanswered",
        "contents": "用于教师收件箱未回答状态的冒烟提问。",
        "created_at": "2026-01-02 09:00:00",
        "views": 12,
        "reply_cnt": 0,
        "pinned": False,
        "deleted_at": None,
        "answer": None,
    },
    {
        "title": "[smoke] teacher inbox answered",
        "contents": "用于教师收件箱已回答状态的冒烟提问。",
        "created_at": "2026-01-02 09:05:00",
        "views": 24,
        "reply_cnt": 1,
        "pinned": False,
        "deleted_at": None,
        "answer": "[smoke] teacher answered this question for regression coverage.",
    },
    {
        "title": "[smoke] teacher inbox pinned",
        "contents": "用于教师收件箱已置顶状态的冒烟提问。",
        "created_at": "2026-01-02 09:10:00",
        "views": 36,
        "reply_cnt": 0,
        "pinned": True,
        "deleted_at": None,
        "answer": None,
    },
    {
        "title": "[smoke] teacher inbox deleted",
        "contents": "用于教师收件箱已删除状态的冒烟提问。",
        "created_at": "2026-01-02 09:15:00",
        "views": 48,
        "reply_cnt": 0,
        "pinned": False,
        "deleted_at": "2026-01-02 10:00:00",
        "answer": None,
    },
]


def ensure_user(
    db_path: Path,
    name: str,
    email: str,
    password: str,
    role: str = "student",
    nickname: str | None = None,
) -> int:
    """幂等地插入/更新用户，返回其 id。"""
    salt, pw_hash = hash_password(password)
    nickname = nickname or name

    conn = sqlite3.connect(str(db_path))
    try:
        conn.execute("PRAGMA foreign_keys = ON")
        row = conn.execute(
            "SELECT id FROM users WHERE name = ? OR email = ?",
            (name, email),
        ).fetchone()

        if row:
            user_id = row[0]
            conn.execute(
                """
                UPDATE users
                SET email = ?, salt = ?, password_hash = ?, role = ?, nickname = ?
                WHERE id = ?
                """,
                (email, salt, pw_hash, role, nickname, user_id),
            )
        else:
            cur = conn.execute(
                """
                INSERT INTO users (name, email, salt, password_hash, role, nickname, introduction)
                VALUES (?, ?, ?, ?, ?, ?, '冒烟测试账户')
                """,
                (name, email, salt, pw_hash, role, nickname),
            )
            user_id = cur.lastrowid

        # settings 必须同步存在，否则 /user 会拿不到
        exists = conn.execute(
            "SELECT 1 FROM settings WHERE id = ?", (user_id,)
        ).fetchone()
        if not exists:
            conn.execute(
                """
                INSERT INTO settings (id, theme_id, notify_switch, notify_email)
                VALUES (?, 1, 1, ?)
                """,
                (user_id, email),
            )

        conn.commit()
        return int(user_id)
    finally:
        conn.close()


def ensure_teacher(
    db_path: Path,
    name: str,
    email: str,
    password: str,
) -> int:
    """幂等地创建一个 teacher 用户 + teachers 行，返回其 id。"""
    uid = ensure_user(db_path, name, email, password, role="teacher")

    conn = sqlite3.connect(str(db_path))
    try:
        conn.execute("PRAGMA foreign_keys = ON")
        exists = conn.execute(
            "SELECT 1 FROM teachers WHERE id = ?", (uid,)
        ).fetchone()
        if not exists:
            conn.execute(
                """
                INSERT INTO teachers (id, responses, perm)
                VALUES (?, 0, 'public')
                """,
                (uid,),
            )
        else:
            conn.execute(
                """
                UPDATE teachers
                SET perm = 'public'
                WHERE id = ?
                """,
                (uid,),
            )
        conn.commit()
        return uid
    finally:
        conn.close()


def ensure_teacher_questions(db_path: Path, student_id: int, teacher_id: int) -> None:
    """为 smoke teacher 幂等写入覆盖各收件箱状态的冒烟提问。"""
    conn = sqlite3.connect(str(db_path))
    try:
        conn.execute("PRAGMA foreign_keys = ON")
        for item in SMOKE_QUESTION_SEEDS:
            row = conn.execute(
                "SELECT id FROM questions WHERE title = ?",
                (item["title"],),
            ).fetchone()
            values = (
                student_id,
                teacher_id,
                item["contents"],
                item["created_at"],
                item["views"],
                item["reply_cnt"],
                item["deleted_at"],
            )
            if row:
                question_id = int(row[0])
                conn.execute(
                    """
                    UPDATE questions
                    SET src_user_id = ?, dst_user_id = ?, contents = ?, created_at = ?,
                        views = ?, reply_cnt = ?, deleted_at = ?
                    WHERE id = ?
                    """,
                    (*values, question_id),
                )
            else:
                cur = conn.execute(
                    """
                    INSERT INTO questions
                        (src_user_id, dst_user_id, title, contents, created_at, views, reply_cnt, deleted_at)
                    VALUES (?, ?, ?, ?, ?, ?, ?, ?)
                    """,
                    (
                        student_id,
                        teacher_id,
                        item["title"],
                        item["contents"],
                        item["created_at"],
                        item["views"],
                        item["reply_cnt"],
                        item["deleted_at"],
                    ),
                )
                question_id = int(cur.lastrowid)

            conn.execute(
                "DELETE FROM favorites WHERE user_id = ? AND question_id = ? AND package = 'top'",
                (teacher_id, question_id),
            )
            if item["pinned"]:
                conn.execute(
                    """
                    INSERT INTO favorites (user_id, question_id, created_at, package)
                    VALUES (?, ?, ?, 'top')
                    ON CONFLICT(user_id, question_id, package) DO UPDATE SET created_at = excluded.created_at
                    """,
                    (teacher_id, question_id, item["created_at"]),
                )

            conn.execute(
                "DELETE FROM answers WHERE question_id = ? AND contents LIKE '[smoke]%'",
                (question_id,),
            )
            if item["answer"]:
                conn.execute(
                    """
                    INSERT INTO answers (user_id, question_id, contents, created_at, upvotes)
                    VALUES (?, ?, ?, ?, 0)
                    """,
                    (teacher_id, question_id, item["answer"], item["created_at"]),
                )

        conn.commit()
    finally:
        conn.close()
