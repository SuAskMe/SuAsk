# SQLite 迁移脚本

这里放"一次性"的 schema 变更，按编号顺序跑。

## 命名规范

`NNN_description.sql` 或 `NNN_description.py`（需要 Python 辅助的场景，比如存档数据、备份 db 文件）。

## 当前清单

| 编号 | 说明 | 幂等？ |
| --- | --- | --- |
| 003 | 移除"问大家"模块 + 为删除业务和公告模块铺路（FK CASCADE/SET NULL、`deleted_at`、`announcement_id`、索引重整） | ❌ 只能跑一次 |

`001` 和 `002` 当前不存在，保留给后续变更占位。

## 怎么跑（Windows PowerShell）

```powershell
# 推荐：用 Python runner 封装备份 + 存档 + 迁移
python database/sqlite/migrations/003_drop_public_and_prepare.py `
    --db ./database/suask.db `
    --archive ./database/archive

# 跑完看到 [OK] 就是成功；suask.db 已原地升级到新 schema
# 如果中途失败：原 db 不会被破坏（先复制再改）；备份文件放在和 db 同目录
```

### 不想动 python 的情况

也可以直接把 `003_*.sql` 过 sqlite3 CLI：

```powershell
Copy-Item ./database/suask.db ./database/suask.db.bak
sqlite3 ./database/suask.db ".read ./database/sqlite/migrations/003_drop_public_and_prepare.sql"
sqlite3 ./database/suask.db "PRAGMA foreign_key_check;"   # 必须输出 0 行
```

注意纯 SQL 方式**不会把"问大家"数据存成 JSON**，只是硬删。如果要存档请走 Python runner。
