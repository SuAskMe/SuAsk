"""
Migration 007: 迁移教师头像
- 从 teachers.avatar_url 下载图片
- 压缩后存入 upload/ 目录（按 hash 分目录）
- 写入 files 表
- 更新 users.avatar_file_id
- 最后删除 teachers.avatar_url 列

运行方式: python database/sqlite/migrations/007_migrate_teacher_avatars.py
需要网络连接（下载图片）和 Pillow 库（压缩图片）

依赖: pip install Pillow requests
"""

import hashlib
import os
import sqlite3
import sys
from pathlib import Path

try:
    import requests
    from PIL import Image
    from io import BytesIO
except ImportError:
    print("请先安装依赖: pip install Pillow requests")
    sys.exit(1)

# 配置
DB_PATH = "suask.db"
UPLOAD_DIR = "upload"
MAX_SIZE = (400, 400)  # 压缩后最大尺寸
JPEG_QUALITY = 80      # JPEG 压缩质量

def compute_hash(data: bytes) -> bytes:
    """计算文件 SHA-256 hash（与项目 utility/files 一致）"""
    return hashlib.sha256(data).digest()

def hash_to_string(h: bytes) -> str:
    return h.hex()

def get_upload_path(hash_bytes: bytes, filename: str) -> str:
    """按项目规则生成存储路径: upload/xx/xx/filename"""
    hs = hash_to_string(hash_bytes)
    return os.path.join(UPLOAD_DIR, hs[0:2], hs[2:4])

def compress_image(image_data: bytes, max_size=MAX_SIZE, quality=JPEG_QUALITY) -> bytes:
    """压缩图片到指定尺寸和质量"""
    img = Image.open(BytesIO(image_data))
    # 转换为 RGB（去掉 alpha 通道）
    if img.mode in ('RGBA', 'P'):
        img = img.convert('RGB')
    # 等比缩放
    img.thumbnail(max_size, Image.Resampling.LANCZOS)
    # 输出为 JPEG
    output = BytesIO()
    img.save(output, format='JPEG', quality=quality, optimize=True)
    return output.getvalue()

def download_image(url: str) -> bytes | None:
    """下载图片，返回原始字节"""
    try:
        resp = requests.get(url, timeout=30)
        if resp.status_code == 200:
            return resp.content
        print(f"  下载失败 HTTP {resp.status_code}: {url}")
        return None
    except Exception as e:
        print(f"  下载异常: {url} -> {e}")
        return None

def rename_file(hash_bytes: bytes, original_name: str) -> str:
    """生成新文件名: hash前16字符 + 原扩展名"""
    hs = hash_to_string(hash_bytes)
    ext = ".jpg"  # 压缩后统一为 jpg
    return hs[:16] + ext

def main():
    if not os.path.exists(DB_PATH):
        print(f"数据库文件不存在: {DB_PATH}")
        sys.exit(1)

    conn = sqlite3.connect(DB_PATH)
    cur = conn.cursor()

    # 获取所有有 avatar_url 的教师
    cur.execute("SELECT id, avatar_url FROM teachers WHERE avatar_url IS NOT NULL AND avatar_url != ''")
    teachers = cur.fetchall()
    print(f"找到 {len(teachers)} 位教师需要迁移头像")

    success_count = 0
    skip_count = 0
    fail_count = 0

    for teacher_id, avatar_url in teachers:
        # 检查是否已经有 avatar_file_id
        cur.execute("SELECT avatar_file_id FROM users WHERE id = ?", (teacher_id,))
        row = cur.fetchone()
        if row and row[0]:
            print(f"  [跳过] id={teacher_id} 已有 avatar_file_id={row[0]}")
            skip_count += 1
            continue

        print(f"  [处理] id={teacher_id}, url={avatar_url}")

        # 下载图片
        image_data = download_image(avatar_url)
        if not image_data:
            fail_count += 1
            continue

        # 压缩图片
        try:
            compressed = compress_image(image_data)
        except Exception as e:
            print(f"  压缩失败: {e}")
            fail_count += 1
            continue

        # 计算 hash
        file_hash = compute_hash(compressed)
        new_filename = rename_file(file_hash, avatar_url.split('/')[-1])

        # 存储文件
        upload_path = get_upload_path(file_hash, new_filename)
        os.makedirs(upload_path, exist_ok=True)
        file_path = os.path.join(upload_path, new_filename)
        with open(file_path, 'wb') as f:
            f.write(compressed)

        # 获取原始文件名（用于 files 表）
        original_name = avatar_url.split('/')[-1]

        # 写入 files 表
        cur.execute(
            "INSERT INTO files (name, hash, uploader_id, created_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP)",
            (original_name, file_hash, teacher_id)
        )
        file_id = cur.lastrowid

        # 更新 users.avatar_file_id
        cur.execute(
            "UPDATE users SET avatar_file_id = ? WHERE id = ?",
            (file_id, teacher_id)
        )

        print(f"    -> file_id={file_id}, path={file_path}, size={len(compressed)} bytes")
        success_count += 1

    conn.commit()

    # 打印统计
    print(f"\n=== 迁移完成 ===")
    print(f"成功: {success_count}")
    print(f"跳过: {skip_count}")
    print(f"失败: {fail_count}")

    # 询问是否删除 avatar_url 列
    if success_count + skip_count == len(teachers) and fail_count == 0:
        print("\n所有教师头像已迁移，可以安全删除 teachers.avatar_url 列。")
        print("执行以下 SQL 删除列（SQLite 需要重建表）：")
        print("""
-- SQLite 不支持 DROP COLUMN（3.35.0+ 支持），如果版本够新：
ALTER TABLE teachers DROP COLUMN avatar_url;

-- 如果版本不够新，需要重建表：
-- CREATE TABLE teachers_new (id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE, perm TEXT, responses INTEGER DEFAULT 0);
-- INSERT INTO teachers_new SELECT id, perm, responses FROM teachers;
-- DROP TABLE teachers;
-- ALTER TABLE teachers_new RENAME TO teachers;
""")
    else:
        print(f"\n有 {fail_count} 个失败，请检查后重新运行。")

    conn.close()

if __name__ == "__main__":
    main()
