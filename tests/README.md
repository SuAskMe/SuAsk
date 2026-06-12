# SuAsk 测试

两层防护网，在重构前后都应该跑一遍。

## 1. Go 单元测试（快，零依赖）

覆盖几个马上要改动的纯函数工具：密码哈希、分页计算、字符串截断、文件 hash。

```powershell
go test ./utility/... ./utility/files/...
```

## 2. HTTP 冒烟回归（需要后端 + Redis + SQLite 运行）

针对跑起来的后端做端到端检查。使用方式见 [`smoke/README.md`](./smoke/README.md)。

典型工作流：

```powershell
# 改代码前，先录制一次基线
python tests/smoke/run_smoke.py seed
python tests/smoke/run_smoke.py snapshot

# 改完后验证结构没有变
python tests/smoke/run_smoke.py verify

# 静态文件暴露面回归检查（独立运行，不依赖 snapshot）
python tests/smoke/security_check.py
```

## 注意

完整 Go 测试基线应保持可直接运行：

```powershell
go test ./...
```

如果本地没有 `manifest/config/config.yaml`，可以显式使用脱敏模板：

```powershell
$env:GF_GCFG_FILE="manifest/config/config.yaml.example"
go test ./...
Remove-Item Env:GF_GCFG_FILE
```
