"""
直接解析 mysqldump 生成的 dump.sql，灌入 SQLite（无需启动 MySQL 实例）。

用法：
    python migrate_from_dump.py --dump ../../dump.sql --sqlite ../suask.db --overwrite

实现要点：
- dump 文件是二进制安全的：`_binary '...'` 字面量里会出现任意字节（BLAKE2b hash 等），
  因此整个解析过程都以 bytes 为单位进行，再在最后把 TEXT 列按 utf-8 解码。
- 仅支持 `mysqldump --extended-insert`（默认打开）生成的 `VALUES (...), (...)` 格式。
- bit(1) 字段会被自动判别为 0 / 1。
- 旧数据归一化：
    * notifications.reply_to_id / answer_id 中的 0 → NULL
    * favorites.package 的 '置顶' → 'top'
"""

from __future__ import annotations

import argparse
import re
import sqlite3
import sys
from pathlib import Path


# ---------------- 字段元信息 ---------------- #

BIT_FIELDS = {
    ("questions", "is_private"),
    ("notifications", "is_read"),
    ("settings", "notify_switch"),
    ("config", "id"),
}

ZERO_AS_NULL = {
    ("notifications", "reply_to_id"),
    ("notifications", "answer_id"),
}

VALUE_REWRITES = {
    ("favorites", "package"): {"置顶": "top"},
}

# 按 FK 依赖顺序灌表
TABLE_ORDER = [
    "users", "files", "teachers", "settings", "themes", "config",
    "questions", "answers", "attachments", "favorites",
    "upvotes", "notifications", "user_relation", "viewed",
]

# BLOB 列（不要尝试按 utf-8 解码）
BLOB_FIELDS = {
    ("files", "hash"),
}


# ---------------- bytes 解析器 ---------------- #

_MYSQL_ESC = {
    ord("\\"): b"\\",
    ord("'"): b"'",
    ord('"'): b'"',
    ord("n"): b"\n",
    ord("r"): b"\r",
    ord("t"): b"\t",
    ord("0"): b"\x00",
    ord("Z"): b"\x1a",
    ord("b"): b"\x08",
}


def _read_quoted(buf: bytes, pos: int) -> tuple[bytes, int]:
    """从 pos 指向的单引号开始读一个 MySQL 字符串，返回 (字节串, 引号后位置)。"""
    assert buf[pos:pos + 1] == b"'", f"字符串必须以单引号开始，实际 {buf[pos:pos+1]!r}"
    i = pos + 1
    out = bytearray()
    n = len(buf)
    while i < n:
        b = buf[i]
        if b == 0x5C and i + 1 < n:        # '\'
            out.extend(_MYSQL_ESC.get(buf[i + 1], bytes([buf[i + 1]])))
            i += 2
            continue
        if b == 0x27:                       # '\''
            if i + 1 < n and buf[i + 1] == 0x27:   # '' → 一个单引号
                out.append(0x27)
                i += 2
                continue
            return bytes(out), i + 1
        out.append(b)
        i += 1
    raise ValueError(f"未闭合的字符串 at pos={pos}")


def _parse_value(buf: bytes, pos: int) -> tuple[object, int]:
    n = len(buf)
    while pos < n and buf[pos] in b" \t\r\n":
        pos += 1
    # NULL
    if buf[pos:pos + 4].upper() == b"NULL":
        return None, pos + 4
    # _binary '...'
    if buf[pos:pos + 7].upper() == b"_BINARY":
        pos += 7
        while pos < n and buf[pos] in b" \t":
            pos += 1
        data, pos = _read_quoted(buf, pos)
        return data, pos
    # 0x...
    if buf[pos:pos + 2].lower() == b"0x":
        j = pos + 2
        while j < n and chr(buf[j]) in "0123456789abcdefABCDEF":
            j += 1
        return bytes.fromhex(buf[pos + 2:j].decode("ascii")), j
    # 字符串
    if buf[pos] == 0x27:
        data, pos = _read_quoted(buf, pos)
        return data, pos
    # 数字
    m = re.match(rb"-?\d+(?:\.\d+)?", buf[pos:])
    if m:
        token = m.group()
        pos += len(token)
        s = token.decode("ascii")
        if b"." in token:
            return float(s), pos
        return int(s), pos
    raise ValueError(f"无法识别字面量 at pos={pos}: {buf[pos:pos+40]!r}")


def _parse_row(buf: bytes, pos: int) -> tuple[list, int]:
    assert buf[pos:pos + 1] == b"("
    pos += 1
    row = []
    while True:
        val, pos = _parse_value(buf, pos)
        row.append(val)
        n = len(buf)
        while pos < n and buf[pos] in b" \t\r\n":
            pos += 1
        if buf[pos:pos + 1] == b",":
            pos += 1
            continue
        if buf[pos:pos + 1] == b")":
            return row, pos + 1
        raise ValueError(f"解析行失败 at pos={pos}: {buf[pos:pos+40]!r}")


_INSERT_RE = re.compile(
    rb"INSERT\s+INTO\s+`(?P<table>\w+)`(?:\s*\([^)]*\))?\s+VALUES\s+(?P<rows>.*?);\s*\n",
    re.IGNORECASE | re.DOTALL,
)


def extract_inserts(dump: bytes) -> dict[str, list[list]]:
    out: dict[str, list[list]] = {t: [] for t in TABLE_ORDER}
    for m in _INSERT_RE.finditer(dump):
        table = m.group("table").decode("ascii")
        if table not in out:
            continue
        values = m.group("rows")
        pos = 0
        n = len(values)
        while pos < n:
            while pos < n and values[pos] in b" \t\r\n,":
                pos += 1
            if pos >= n:
                break
            if values[pos:pos + 1] != b"(":
                raise ValueError(f"{table} VALUES 异常 at {pos}: {values[pos:pos+40]!r}")
            row, pos = _parse_row(values, pos)
            out[table].append(row)
    return out


_CREATE_RE = re.compile(
    rb"CREATE\s+TABLE\s+`(?P<tbl>\w+)`\s*\((?P<body>.*?)\)\s*ENGINE",
    re.DOTALL | re.IGNORECASE,
)


def extract_column_order(dump: bytes) -> dict[str, list[str]]:
    order: dict[str, list[str]] = {}
    for m in _CREATE_RE.finditer(dump):
        tbl = m.group("tbl").decode("ascii")
        cols: list[str] = []
        for raw_line in m.group("body").splitlines():
            line = raw_line.strip().rstrip(b",")
            cm = re.match(rb"`(\w+)`", line)
            upper = line.upper()
            if not cm:
                continue
            if upper.startswith((b"PRIMARY KEY", b"KEY", b"UNIQUE",
                                 b"FULLTEXT", b"CONSTRAINT", b"CHECK")):
                continue
            cols.append(cm.group(1).decode("ascii"))
        order[tbl] = cols
    return order


# ---------------- 值归一化 ---------------- #

def normalize(table: str, col: str, v):
    # bit(1) → 0/1
    if (table, col) in BIT_FIELDS and isinstance(v, (bytes, bytearray)):
        return 0 if len(v) == 0 else int(v[0])
    # 哨兵 0 → NULL
    if (table, col) in ZERO_AS_NULL and v == 0:
        return None
    # BLOB 列保持原字节
    if (table, col) in BLOB_FIELDS:
        return bytes(v) if isinstance(v, (bytes, bytearray)) else v
    # 其它 bytes（TEXT 列）→ 按 utf-8 解码
    if isinstance(v, (bytes, bytearray)):
        try:
            v = v.decode("utf-8")
        except UnicodeDecodeError:
            v = v.decode("utf-8", errors="replace")
    # 值重写
    rw = VALUE_REWRITES.get((table, col))
    if rw and v in rw:
        return rw[v]
    return v


# ---------------- 主流程 ---------------- #

def main() -> None:
    p = argparse.ArgumentParser()
    p.add_argument("--dump", required=True, help="mysqldump 生成的 .sql")
    p.add_argument("--sqlite", default=str(Path(__file__).parent.parent / "suask.db"))
    p.add_argument("--schema", default=str(Path(__file__).parent / "schema.sqlite.sql"))
    p.add_argument("--overwrite", action="store_true")
    args = p.parse_args()

    dump_path = Path(args.dump).resolve()
    sqlite_path = Path(args.sqlite).resolve()
    schema_path = Path(args.schema).resolve()

    if not dump_path.exists():
        sys.exit(f"dump 文件不存在：{dump_path}")
    if sqlite_path.exists():
        if args.overwrite:
            sqlite_path.unlink()
        else:
            sys.exit(f"{sqlite_path} 已存在，加 --overwrite 覆盖")

    print(f"[*] 读取 dump：{dump_path}")
    dump_bytes = dump_path.read_bytes()

    columns = extract_column_order(dump_bytes)
    rows_by_table = extract_inserts(dump_bytes)

    print(f"[*] 创建 SQLite：{sqlite_path}")
    dst = sqlite3.connect(sqlite_path)
    dst.execute("PRAGMA foreign_keys = OFF")
    dst.executescript(schema_path.read_text(encoding="utf-8"))
    dst.execute("PRAGMA foreign_keys = OFF")

    total = 0
    for table in TABLE_ORDER:
        rows = rows_by_table.get(table) or []
        if not rows:
            print(f"[.] {table}: 0 行")
            continue
        cols = columns.get(table)
        if not cols:
            print(f"[!] {table} 没找到字段顺序，跳过")
            continue
        col_list = ",".join(f'"{c}"' for c in cols)
        placeholders = ",".join("?" * len(cols))
        # user_relation 在 MySQL 原表没有主键，可能有重复行；SQLite 加了 PK，这里用 IGNORE 去重
        verb = "INSERT OR IGNORE" if table == "user_relation" else "INSERT"
        sql = f'{verb} INTO "{table}" ({col_list}) VALUES ({placeholders})'
        data = [
            tuple(normalize(table, cols[i], v) for i, v in enumerate(row))
            for row in rows
        ]
        dst.executemany(sql, data)
        # 同步 AUTOINCREMENT 的 sqlite_sequence
        if "id" in cols:
            idx = cols.index("id")
            int_ids = [r[idx] for r in rows if isinstance(r[idx], int)]
            if int_ids:
                dst.execute(
                    "INSERT OR REPLACE INTO sqlite_sequence (name, seq) VALUES (?, ?)",
                    (table, max(int_ids)),
                )
        print(f"[+] {table}: {len(rows)} 行")
        total += len(rows)

    dst.commit()
    dst.execute("PRAGMA foreign_keys = ON")
    problems = dst.execute("PRAGMA foreign_key_check").fetchall()
    if problems:
        print("\n[!] 外键检查发现问题（前 20 条）：")
        for row in problems[:20]:
            print("   ", row)
    else:
        print("\n[OK] 外键完整性检查通过")
    dst.close()
    print(f"\n[DONE] 共迁移 {total} 行 → {sqlite_path}")


if __name__ == "__main__":
    main()
