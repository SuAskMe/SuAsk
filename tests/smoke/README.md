# SuAsk 后端回归冒烟测试

这套脚本是重构时的"安全网"：在动代码之前，先把当前接口的行为录制成 snapshot，
每次改动后再跑一遍，自动告诉你哪个接口的返回结构变了。

## 运行前提

1. 后端已启动（默认 `http://127.0.0.1:8080`）；
2. Redis 已启动（登录流程依赖）；
3. SQLite 数据库文件存在（默认 `./database/suask.db`）；
4. Python 3.9+，安装依赖：

```powershell
pip install -r tests/smoke/requirements.txt
```

## 用法

脚本从**项目根目录**运行（默认 cwd）。

### 1. 第一次运行 —— 录制快照（在你改代码之前）

```powershell
# 先在 SQLite 里塞一个确定性的测试账户（幂等，多跑也没关系）
python tests/smoke/run_smoke.py seed

# 录制快照
python tests/smoke/run_smoke.py snapshot
```

会把每个接口的响应"形状"（字段名+类型+必要的结构）写到
`tests/smoke/snapshots/*.json`，值不记录，避免时间戳/ID 之类的噪声。

### 2. 改完代码后 —— 验证

```powershell
python tests/smoke/run_smoke.py verify
```

任何结构变更（字段消失/类型变化/新增未预期字段）都会打印 diff 并以非零退出。

### 3. 参数

```powershell
python tests/smoke/run_smoke.py verify `
    --base-url http://127.0.0.1:8080 `
    --db ./database/suask.db `
    --user smoke_student `
    --password Smoke-Pass-123!
```

默认值都在脚本顶部。

## 覆盖的接口

### 读
- `POST /login`
- `GET /user`
- `GET /info/user?id=...`
- `GET /info/teacher`
- `GET /questions/public?sort_type=0&page=1`
- `GET /answer?question_id=...`
- `GET /favorites?sort_type=0&page=1`
- `GET /history?sort_type=0&page=1`
- `GET /notification?user_id=...`
- `GET /notification/count?user_id=...`

### 写（幂等：同一操作执行两次恢复原状）
- `POST /favorites` → 收藏/取消收藏
- `POST /answer/upvote` → 点赞/取消点赞
- `POST /questions/add` → 新发一条问老师的问题（之后不删，但可识别）

## 静态文件暴露面检查

独立于 snapshot 机制，单独一个脚本：

```powershell
python tests/smoke/security_check.py
```

会验证 `/database/suask.db`、`/manifest/config/config.yaml`、`/main.exe`、`/` 等路径
都返回 4xx，同时确认 `/upload/...` 下的已有文件仍可正常访问。一旦有人不小心又把
`SetServerRoot(".")` 改回来，这里会立刻红。

## 不测的东西

- 邮件验证码注册 / 改密码（依赖外部 SMTP）
- 文件上传的图片内容（仅测 multipart 通路，不校验图片本身）
- 高并发场景
