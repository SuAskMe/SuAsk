"""
给回归测试准备一个确定的用户账号，直接写 SQLite（绕过邮件验证码）。

用户名/密码默认是 smoke_student / Smoke-Pass-123! ，可用 CLI 参数覆盖。

密码哈希算法必须和 utility.EncryptPassword 保持一致：
  md5(md5(password) + md5(salt))
一旦后端改成 bcrypt，这里也要同步更新，否则登录会失败 —— 这就是我们要的"提醒"。
"""

from __future__ import annotations

import hashlib
import sqlite3
from pathlib import Path


def md5_hex(s: str) -> str:
    return hashlib.md5(s.encode("utf-8")).hexdigest()


def encrypt_password(password: str, salt: str) -> str:
    """必须与 utility/utils.go 的 EncryptPassword 完全一致"""
    return md5_hex(md5_hex(password) + md5_hex(salt))


def ensure_user(
    db_path: Path,
    name: str,
    email: str,
    password: str,
    role: str = "student",
    nickname: str | None = None,
) -> int:
    """幂等地插入/更新用户，返回其 id。"""
    salt = "smokesalt1"  # 固定 10 位（和 grand.S(10) 保持长度一致）
    pw_hash = encrypt_password(password, salt)
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
