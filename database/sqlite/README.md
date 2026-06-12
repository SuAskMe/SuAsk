# SuAsk SQLite 迁移说明

## 文件说明

| 文件 | 用途 |
| --- | --- |
| `schema.sqlite.sql` | SQLite 建表脚本（已从 MySQL DDL 调整过） |
| `migrate.py`        | 从运行中的 MySQL 抽数据灌入 SQLite 的一次性脚本 |

## 执行步骤

### 1. 从 MySQL 同步数据

前置条件：本地 MySQL 已 import `dump.sql` 或现有生产库可连接。

```powershell
pip install pymysql
python migrate.py --overwrite
```

参数若和默认值不同（默认和 `manifest/config/config.yaml` 的 MySQL 配置一致）：

```powershell
python migrate.py \
  --mysql-host 127.0.0.1 --mysql-port 3306 \
  --mysql-user root --mysql-pass 123456 \
  --mysql-db   suask \
  --sqlite     ../suask.db \
  --schema     schema.sqlite.sql \
  --overwrite
```

结束后会在 `database/suask.db` 生成新的 SQLite 数据库，并打印每张表迁移的行数与外键完整性检查结果。

### 2. 从 dump.sql 直接导入（不经过 MySQL）

如果手上只有 `dump.sql`，**推荐的方式仍然是先导回本地 MySQL 再跑上面的 Python 脚本**。原因：

- `dump.sql` 里的 `bit(1)` 字段用 `_binary ''` / `_binary '\0'` 表示，SQLite 无法直接识别；
- `binary(32)` 的 blake2b hash 在 dump 里是非打印字符混合（不是 `X'...'`），无法直接 `INSERT`；
- `LOCK/UNLOCK TABLES`、`SET @OLD_...=@@...`、`/*!... */` 等 MySQL 方言也需要过滤。

可以这样导回 MySQL：

```bash
mysql -uroot -p -e "CREATE DATABASE suask CHARACTER SET utf8mb4 COLLATE utf8mb4_zh_0900_as_cs"
mysql -uroot -p suask < dump.sql
```

然后执行上面的 `migrate.py`。

### 3. 应用侧配置

参见 `manifest/config/config.yaml` 里的 `database` 配置项，已切换到：

```yaml
database:
  default:
    type: "sqlite"
    link: "sqlite::@file(./database/suask.db)"
```

并且 `main.go` 已经把 `mysql` 驱动改成了 `sqlite` 驱动。

## schema 关键差异速查

| MySQL                                   | SQLite                                            |
| --------------------------------------- | ------------------------------------------------- |
| `int AUTO_INCREMENT`                    | `INTEGER PRIMARY KEY AUTOINCREMENT`               |
| `bit(1)` + `_binary ''`                 | `INTEGER CHECK (col IN (0,1))`                    |
| `binary(32)`                            | `BLOB`                                            |
| `enum('a','b')`                         | `TEXT CHECK (col IN ('a','b'))`                   |
| `varchar(N)` / `text`                   | `TEXT`                                            |
| `FULLTEXT KEY … WITH PARSER ngram`      | 取消，代码侧改为 `LIKE '%keyword%'`                 |
| `((a IS NOT NULL) + (b IS NOT NULL)=1)` | `((a IS NULL) <> (b IS NULL))`                    |
| `ON DELETE RESTRICT` 默认生效             | 需连接时 `PRAGMA foreign_keys=ON`                  |

## 已知数据归一

- `questions.is_private`：当前 SQLite schema 已彻底移除该历史字段；从旧 MySQL / dump 导入时会自动忽略。
- `notifications.reply_to_id / answer_id`：旧数据出现了 `0` 作为"无"的哨兵，迁移时转 `NULL`。
- `favorites.package`：旧数据里混用中文 `'置顶'` 和英文 `'default'`。迁移脚本把 `'置顶' → 'top'`，并在 schema 里允许两种过渡值。代码侧建议逐步改用 `'top'`。
