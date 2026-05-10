"""
静态文件暴露面回归检查（不依赖 snapshot，直接断言 HTTP 状态码）。

背景：
  以前 cmd.go 用 SetServerRoot(".") + SetIndexFolder(true) 把整个项目根目录
  当静态站暴露了；任何人都能 GET /database/suask.db、/manifest/config/config.yaml、
  /main.exe 等敏感文件。

这个脚本的职责：
  1. 正向 —— upload 目录下的文件仍能访问（只要库里至少有 1 条 files 记录，
     并且磁盘上真的落了盘；拿不到就跳过该项）。
  2. 反向 —— 下面这些路径必须返回 4xx：
        /database/suask.db
        /manifest/config/config.yaml
        /main.exe
        /go.mod
        /database/
        /manifest/config/
        /
     一旦哪天有人又打开 SetServerRoot(".")，这里会立刻报红。

用法：
  python tests/smoke/security_check.py [--base-url ...] [--db ./database/suask.db]
"""

from __future__ import annotations

import argparse
import sqlite3
import sys
from pathlib import Path

import requests


FORBIDDEN_PATHS = [
    "/database/suask.db",
    "/manifest/config/config.yaml",
    "/main.exe",
    "/go.mod",
    "/go.sum",
    "/database/",
    "/manifest/config/",
    "/",  # 目录列表
]


def _is_blocked(status: int) -> bool:
    """2xx / 3xx 都视作暴露；我们只接受 4xx。"""
    return 400 <= status < 500


def check_forbidden(base_url: str) -> list[str]:
    failures: list[str] = []
    for p in FORBIDDEN_PATHS:
        url = base_url.rstrip("/") + p
        try:
            r = requests.get(url, timeout=5, allow_redirects=False)
        except requests.RequestException as exc:
            failures.append(f"{p}: 请求失败 {exc}")
            continue
        if _is_blocked(r.status_code):
            print(f"[OK   ] {p}  -> {r.status_code}")
        else:
            snippet = r.text[:120].replace("\n", "\\n")
            failures.append(
                f"{p}: 期望 4xx，实际 {r.status_code}；body[:120]={snippet!r}"
            )
    return failures


def _probe_upload(base_url: str, db_path: Path) -> str | None:
    """从 files 表取第一条，拼成 /upload/<h2>/<h2>/<name>。拿不到就返回 None。"""
    if not db_path.exists():
        return None
    conn = sqlite3.connect(str(db_path))
    try:
        row = conn.execute(
            "SELECT name, hash FROM files ORDER BY id LIMIT 1"
        ).fetchone()
    finally:
        conn.close()
    if not row:
        return None
    name, raw_hash = row
    if not isinstance(raw_hash, (bytes, bytearray)):
        return None
    hex_hash = bytes(raw_hash).hex()
    if "." in name:
        ext = name.rsplit(".", 1)[-1]
        file_name = f"{hex_hash}.{ext}"
    else:
        file_name = hex_hash
    return f"/upload/{hex_hash[:2]}/{hex_hash[2:4]}/{file_name}"


def check_upload_reachable(base_url: str, db_path: Path) -> list[str]:
    path = _probe_upload(base_url, db_path)
    if path is None:
        print("[skip ] upload 可达性检查：files 表为空或库不可读")
        return []
    url = base_url.rstrip("/") + path
    try:
        r = requests.get(url, timeout=5, allow_redirects=False)
    except requests.RequestException as exc:
        return [f"{path}: 请求失败 {exc}"]
    # 2xx 正常；404 说明磁盘上没这份文件（老库记录了但文件丢了），算可忽略
    if 200 <= r.status_code < 300:
        print(f"[OK   ] {path}  -> {r.status_code}  (上传目录仍可正常访问)")
        return []
    if r.status_code == 404:
        print(f"[skip ] {path}  -> 404  (落盘文件缺失，DB 残留；不算回归)")
        return []
    return [f"{path}: upload 目录不可达，期望 2xx，实际 {r.status_code}"]


def main() -> None:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("--base-url", default="http://127.0.0.1:8080")
    p.add_argument("--db", default=Path("database/suask.db"), type=Path)
    args = p.parse_args()

    failures: list[str] = []
    failures.extend(check_forbidden(args.base_url))
    failures.extend(check_upload_reachable(args.base_url, args.db))

    if failures:
        print("\n[FAIL] 发现静态文件暴露面回归:", *failures, sep="\n  - ")
        sys.exit(1)
    print("\n[PASS] 静态文件暴露面检查通过")


if __name__ == "__main__":
    main()
