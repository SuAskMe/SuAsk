"""
Migration 003 runner
  - 备份 suask.db
  - 把 "问大家" 数据 (dst_user_id IS NULL) 导出为 JSON 存档
  - 执行 003_drop_public_and_prepare.sql
  - 跑 PRAGMA foreign_key_check 验证
  - 跑快速 smoke（表是否还在、row count 是否符合预期）

用法：
    python database/sqlite/migrations/003_drop_public_and_prepare.py \
        --db ./database/suask.db \
        --archive ./database/archive

注意：
  - 默认会拒绝对已经迁移过的 db 再跑一次（通过探测 questions.deleted_at 列是否存在）
  - 加 --force 可强行重跑，但那会把 "问大家" 存档覆盖掉
"""

from __future__ import annotations

import argparse
import datetime as dt
import json
import shutil
import sqlite3
import sys
from pathlib import Path


SCRIPT_DIR = Path(__file__).resolve().parent
DEFAULT_SQL = SCRIPT_DIR / "003_drop_public_and_prepare.sql"


# ----------------------------------------------------------------------------
# helpers
# ----------------------------------------------------------------------------


def column_exists(conn: sqlite3.Connection, table: str, col: str) -> bool:
    rows = conn.execute(f"PRAGMA table_info({table})").fetchall()
    return any(r[1] == col for r in rows)


def count(conn: sqlite3.Connection, table: str, where: str = "") -> int:
    sql = f"SELECT COUNT(*) FROM {table}"
    if where:
        sql += f" WHERE {where}"
    return conn.execute(sql).fetchone()[0]


def backup_db(db_path: Path) -> Path:
    ts = dt.datetime.now().strftime("%Y%m%d_%H%M%S")
    bak = db_path.with_name(f"{db_path.name}.bak.{ts}")
    shutil.copy2(db_path, bak)
    return bak


def archive_public_questions(conn: sqlite3.Connection, archive_dir: Path) -> Path:
    archive_dir.mkdir(parents=True, exist_ok=True)
    ts = dt.datetime.now().strftime("%Y%m%d_%H%M%S")
    out = archive_dir / f"dst_null_questions_{ts}.json"

    # 导出问题 + 关联的 answers / attachments / upvotes / favorites / notifications
    # 数据量小（5 条），直接全量 load 到内存即可
    q_rows = conn.execute(
        """
        SELECT id, src_user_id, dst_user_id, title, contents, is_private,
               created_at, views, reply_cnt
          FROM questions WHERE dst_user_id IS NULL
        """
    ).fetchall()
    qids = [r[0] for r in q_rows]

    questions = [dict(zip(
        ["id", "src_user_id", "dst_user_id", "title", "contents",
         "is_private", "created_at", "views", "reply_cnt"], r)) for r in q_rows]

    def fetch_by_qid(table: str, col: str = "question_id") -> list[dict]:
        if not qids:
            return []
        placeholders = ",".join("?" * len(qids))
        cur = conn.execute(
            f"SELECT * FROM {table} WHERE {col} IN ({placeholders})", qids
        )
        cols = [d[0] for d in cur.description]
        return [dict(zip(cols, row)) for row in cur.fetchall()]

    answers = fetch_by_qid("answers")
    attachments_q = fetch_by_qid("attachments")
    favorites = fetch_by_qid("favorites")
    notifications = fetch_by_qid("notifications")
    upvotes_q = fetch_by_qid("upvotes")

    # 对 answer 级联的附件/upvotes/notifications 也顺手一起存
    aids = [a["id"] for a in answers]

    def fetch_by_aid(table: str, col: str) -> list[dict]:
        if not aids:
            return []
        placeholders = ",".join("?" * len(aids))
        cur = conn.execute(
            f"SELECT * FROM {table} WHERE {col} IN ({placeholders})", aids
        )
        cols = [d[0] for d in cur.description]
        return [dict(zip(cols, row)) for row in cur.fetchall()]

    attachments_a = fetch_by_aid("attachments", "answer_id")
    upvotes_a = fetch_by_aid("upvotes", "answer_id")
    notifications_a = fetch_by_aid("notifications", "answer_id")

    payload = {
        "exported_at": dt.datetime.now().isoformat(timespec="seconds"),
        "questions": questions,
        "answers": answers,
        "attachments_of_questions": attachments_q,
        "attachments_of_answers": attachments_a,
        "favorites": favorites,
        "upvotes_of_questions": upvotes_q,
        "upvotes_of_answers": upvotes_a,
        "notifications_of_questions": notifications,
        "notifications_of_answers": notifications_a,
    }

    # sqlite3 会把 BLOB/bytes 读成 bytes，JSON 不能序列化；这里兜底转成 hex
    def _default(o):
        if isinstance(o, (bytes, bytearray)):
            return o.hex()
        raise TypeError(f"not serializable: {type(o).__name__}")

    out.write_text(
        json.dumps(payload, ensure_ascii=False, indent=2, default=_default),
        encoding="utf-8",
    )
    return out


def smoke_after(conn: sqlite3.Connection, q_before: int, a_before: int) -> list[str]:
    """跑几条事实断言，有问题就塞进 failures 列表返回"""
    failures: list[str] = []

    expected_tables = [
        "users", "files", "teachers", "settings", "themes", "config",
        "questions", "answers", "attachments", "favorites", "upvotes",
        "notifications", "user_relation", "viewed",
    ]
    actual = {r[0] for r in conn.execute(
        "SELECT name FROM sqlite_master WHERE type='table'"
    ).fetchall()}
    for t in expected_tables:
        if t not in actual:
            failures.append(f"table missing after migration: {t}")

    # 新增列
    for table, col in [
        ("questions", "deleted_at"),
        ("answers", "deleted_at"),
        ("answers", "announcement_id"),
        ("attachments", "announcement_id"),
        ("files", "deleted_at"),
    ]:
        if not column_exists(conn, table, col):
            failures.append(f"column missing: {table}.{col}")

    # 问大家清空
    leftover = count(conn, "questions", "dst_user_id IS NULL")
    if leftover != 0:
        failures.append(f"expected 0 questions with dst_user_id IS NULL, got {leftover}")

    # row count 没无故丢
    # 注意：questions 和 answers 会减少（"问大家"数据被移走）
    q_after = count(conn, "questions")
    a_after = count(conn, "answers")
    print(f"  questions: {q_before} -> {q_after}")
    print(f"  answers:   {a_before} -> {a_after}")
    if q_after > q_before:
        failures.append("questions 数量反而增加了，异常")
    if a_after > a_before:
        failures.append("answers 数量反而增加了，异常")

    # 外键全部合法
    fk = conn.execute("PRAGMA foreign_key_check").fetchall()
    if fk:
        failures.append(f"foreign_key_check 不干净：{fk[:5]}{'...' if len(fk) > 5 else ''}")
    return failures


# ----------------------------------------------------------------------------
# main
# ----------------------------------------------------------------------------


def main() -> None:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--db", type=Path, default=Path("database/suask.db"))
    p.add_argument("--archive", type=Path, default=Path("database/archive"))
    p.add_argument("--sql", type=Path, default=DEFAULT_SQL)
    p.add_argument("--force", action="store_true",
                   help="即使看起来已迁移过也强制重跑")
    args = p.parse_args()

    if not args.db.exists():
        sys.exit(f"db 不存在: {args.db}")
    if not args.sql.exists():
        sys.exit(f"sql 脚本不存在: {args.sql}")

    conn = sqlite3.connect(str(args.db))
    conn.execute("PRAGMA foreign_keys = OFF")

    # 已迁移检测
    if column_exists(conn, "questions", "deleted_at") and not args.force:
        sys.exit(
            "db 已经迁移过（questions.deleted_at 已存在）；如需重跑请加 --force"
        )

    q_before = count(conn, "questions")
    a_before = count(conn, "answers")
    pub_before = count(conn, "questions", "dst_user_id IS NULL")
    print(f"[*] 迁移前：questions={q_before}, answers={a_before}, 问大家={pub_before}")

    # 1) 备份
    bak = backup_db(args.db)
    print(f"[*] 已备份到 {bak}")

    # 2) 存档
    archive_file = archive_public_questions(conn, args.archive)
    print(f'[*] 已把"问大家"数据存档到 {archive_file}')

    # 3) 执行 SQL（executescript 会自动 COMMIT，用就用）
    sql = args.sql.read_text(encoding="utf-8")
    try:
        conn.executescript(sql)
    except Exception as exc:
        conn.close()
        # 出事回退到备份
        shutil.copy2(bak, args.db)
        sys.exit(f"[FAIL] 迁移失败，已回滚到 {bak}\n原因: {exc}")

    # 4) smoke
    failures = smoke_after(conn, q_before, a_before)
    conn.close()

    if failures:
        # 回滚
        shutil.copy2(bak, args.db)
        print("\n[FAIL] smoke 失败，已回滚到备份：", *failures, sep="\n  - ")
        sys.exit(1)

    print("\n[OK] 迁移完成，外键检查通过")
    print(f"     备份: {bak}")
    print(f"     存档: {archive_file}")


if __name__ == "__main__":
    main()
