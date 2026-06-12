"""
Migration 009: 删除 questions.is_private 历史字段

运行方式:
    python database/sqlite/migrations/009_drop_question_is_private.py --db ./suask.db

默认会:
  - 检查 questions.is_private 是否仍存在
  - 在原地迁移前备份数据库
  - 执行 009_drop_question_is_private.sql
  - 校验 questions 行数不变、目标列已删除、外键完整性正常
"""

from __future__ import annotations

import argparse
import shutil
import sqlite3
import sys
from datetime import datetime
from pathlib import Path


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser()
    p.add_argument("--db", type=Path, default=Path("suask.db"))
    p.add_argument(
        "--sql",
        type=Path,
        default=Path(__file__).with_name("009_drop_question_is_private.sql"),
    )
    return p.parse_args()


def column_exists(conn: sqlite3.Connection, table: str, column: str) -> bool:
    rows = conn.execute(f"PRAGMA table_info({table})").fetchall()
    return any(row[1] == column for row in rows)


def count(conn: sqlite3.Connection, table: str, where: str = "1=1") -> int:
    return conn.execute(f"SELECT COUNT(*) FROM {table} WHERE {where}").fetchone()[0]


def main() -> None:
    args = parse_args()
    db_path = args.db.resolve()
    sql_path = args.sql.resolve()

    if not db_path.exists():
        sys.exit(f"数据库文件不存在: {db_path}")
    if not sql_path.exists():
        sys.exit(f"迁移 SQL 不存在: {sql_path}")

    conn = sqlite3.connect(db_path)
    try:
        if not column_exists(conn, "questions", "is_private"):
            print("[SKIP] questions.is_private 已不存在，无需重复迁移。")
            return
        q_before = count(conn, "questions")
        private_before = count(conn, "questions", "is_private = 1")
    finally:
        conn.close()

    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    backup_path = db_path.with_name(f"{db_path.name}.bak.{timestamp}_drop_is_private")
    shutil.copy2(db_path, backup_path)
    print(f"[*] 已备份数据库 -> {backup_path}")
    print(f"[*] 迁移前 questions={q_before}, legacy_private={private_before}")

    try:
        conn = sqlite3.connect(db_path)
        try:
            conn.executescript(sql_path.read_text(encoding="utf-8"))
            conn.commit()

            if column_exists(conn, "questions", "is_private"):
                raise RuntimeError("迁移后 questions.is_private 仍然存在")

            q_after = count(conn, "questions")
            if q_after != q_before:
                raise RuntimeError(f"questions 行数异常: {q_before} -> {q_after}")

            fk_issues = conn.execute("PRAGMA foreign_key_check").fetchall()
            if fk_issues:
                detail = "\n".join(f"  {row}" for row in fk_issues[:20])
                raise RuntimeError(f"外键检查失败（前 20 条）:\n{detail}")

            print("[OK] questions.is_private 已删除，questions 行数保持不变，外键检查通过。")
        finally:
            conn.close()
    except Exception as exc:
        shutil.copy2(backup_path, db_path)
        sys.exit(f"[FAIL] 迁移失败，已自动恢复备份: {exc}")


if __name__ == "__main__":
    main()
