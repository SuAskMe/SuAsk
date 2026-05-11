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
                INSERT INTO teachers (id, responses, name, perm)
                VALUES (?, 0, ?, 'public')
                """,
                (uid, name),
            )
        conn.commit()
        return uid
    finally:
        conn.close()
