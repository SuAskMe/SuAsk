"""
SuAsk MySQL → SQLite 数据迁移脚本
================================

用法：
    pip install pymysql
    python migrate.py \
        --mysql-host 127.0.0.1 --mysql-port 3306 \
        --mysql-user root       --mysql-pass 123456 \
        --mysql-db   suask      \
        --sqlite     ../suask.db \
        --schema     schema.sqlite.sql

默认值与 manifest/config/config.yaml 中的 mysql 配置保持一致。
"""

from __future__ import annotations

import argparse
import os
import sys
from pathlib import Path

try:
    import pymysql
    from pymysql.cursors import DictCursor
except ImportError:
    sys.exit("请先安装 pymysql：  pip install pymysql")

import sqlite3


# 按外键依赖顺序排列
TABLE_ORDER = [
    "users",
    "files",           # users.avatar_file_id → files.id (表间循环，建表时先跳过 FK)
    "teachers",
    "settings",
    "themes",
    "config",
    "questions",
    "answers",
    "attachments",
    "favorites",
    "upvotes",
    "notifications",
    "user_relation",
    "viewed",
]

# bit(1) 字段 —— MySQL 返回 bytes(b'\x00' / b'\x01')，需要转为 0 / 1
BIT_FIELDS = {
    "questions":     {"is_private"},
    "notifications": {"is_read"},
    "settings":      {"notify_switch"},
    "config":        {"id"},  # 单行配置 id 在 MySQL 里是 bit(1)
}

# 这些列在 MySQL 里 DEFAULT NULL，但旧数据里出现了 0 作为 "无" 的哨兵值 → 归一为 NULL
ZERO_AS_NULL = {
    "notifications": {"reply_to_id", "answer_id"},
}

# 这些列的值需要归一化：'置顶' → 'top'，以配合代码里统一成英文
VALUE_REWRITES = {
    "favorites": {"package": {"置顶": "top"}},
}


def parse_args() -> argparse.Namespace:
    p = argparse.ArgumentParser(description="SuAsk MySQL → SQLite migrator")
    p.add_argument("--mysql-host", default="127.0.0.1")
    p.add_argument("--mysql-port", type=int, default=3306)
    p.add_argument("--mysql-user", default="root")
    p.add_argument("--mysql-pass", default="123456")
    p.add_argument("--mysql-db",   default="suask")
    p.add_argument("--sqlite",     default=str(Path(__file__).parent.parent / "suask.db"))
    p.add_argument("--schema",     default=str(Path(__file__).parent / "schema.sqlite.sql"))
    p.add_argument("--overwrite",  action="store_true", help="若目标 .db 已存在，覆盖删除")
    return p.parse_args()


def normalize_value(table: str, col: str, val):
    """将 MySQL 的 bit / 哨兵 0 / 中文枚举值转成 SQLite 友好值。"""
    # bit(1) → 0/1
    if table in BIT_FIELDS and col in BIT_FIELDS[table] and isinstance(val, (bytes, bytearray)):
        if len(val) == 0:
            return 0
        return int(val[0])
    # 0 → NULL
    if table in ZERO_AS_NULL and col in ZERO_AS_NULL[table] and val == 0:
        return None
    # 值重写（'置顶' → 'top'）
    if table in VALUE_REWRITES and col in VALUE_REWRITES[table]:
        return VALUE_REWRITES[table][col].get(val, val)
    return val


def migrate(args: argparse.Namespace) -> None:
    sqlite_path = Path(args.sqlite).resolve()
    schema_path = Path(args.schema).resolve()
    if not schema_path.exists():
        sys.exit(f"schema 文件不存在：{schema_path}")

    if sqlite_path.exists():
        if args.overwrite:
            sqlite_path.unlink()
        else:
            sys.exit(f"{sqlite_path} 已存在。加 --overwrite 可覆盖。")

    print(f"[*] 连接 MySQL {args.mysql_host}:{args.mysql_port}/{args.mysql_db} ...")
    src = pymysql.connect(
        host=args.mysql_host,
        port=args.mysql_port,
        user=args.mysql_user,
        password=args.mysql_pass,
        database=args.mysql_db,
        charset="utf8mb4",
        cursorclass=DictCursor,
    )

    print(f"[*] 创建 SQLite：{sqlite_path}")
    dst = sqlite3.connect(sqlite_path)
    dst.execute("PRAGMA foreign_keys = OFF")  # 灌数据期间先关外键

    with schema_path.open(encoding="utf-8") as f:
        dst.executescript(f.read())
    # executescript 结束时可能已 ON，这里再次显式关闭，确保 INSERT 顺序无关
    dst.execute("PRAGMA foreign_keys = OFF")

    total = 0
    try:
        with src.cursor() as cur:
            for table in TABLE_ORDER:
                cur.execute(f"SHOW TABLES LIKE '{table}'")
                if not cur.fetchone():
                    print(f"[-] 源库无表 {table}，跳过")
                    continue
                cur.execute(f"SELECT * FROM `{table}`")
                rows = cur.fetchall()
                if not rows:
                    print(f"[.] {table}: 0 行")
                    continue

                cols = list(rows[0].keys())
                col_list = ",".join(f'"{c}"' for c in cols)
                placeholders = ",".join("?" * len(cols))
                sql = f'INSERT INTO "{table}" ({col_list}) VALUES ({placeholders})'

                data = [
                    tuple(normalize_value(table, c, r[c]) for c in cols)
                    for r in rows
                ]
                dst.executemany(sql, data)
                # 同步自增序列
                pk_ints = [r[cols[0]] for r in rows
                           if isinstance(r[cols[0]], int) and cols[0] == "id"]
                if pk_ints:
                    max_id = max(pk_ints)
                    dst.execute(
                        "INSERT OR REPLACE INTO sqlite_sequence (name, seq) VALUES (?, ?)",
                        (table, max_id),
                    )
                print(f"[+] {table}: {len(rows)} 行")
                total += len(rows)

        dst.commit()
        # 灌完再开外键并做一次完整性检查
        dst.execute("PRAGMA foreign_keys = ON")
        problems = dst.execute("PRAGMA foreign_key_check").fetchall()
        if problems:
            print("\n[!] 外键检查发现以下问题（前 20 条）：")
            for row in problems[:20]:
                print("   ", row)
        else:
            print("\n[OK] 外键完整性检查通过")

        print(f"\n[DONE] 共迁移 {total} 行 → {sqlite_path}")
    finally:
        dst.close()
        src.close()


if __name__ == "__main__":
    migrate(parse_args())
