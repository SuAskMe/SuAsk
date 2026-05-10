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

## 注意：仓库已有两个坏掉的测试

- `utility/trie_mux/mux_test.go::TestTrieMux_HasPrefix` —— `HasPrefix("/")`
  的语义在实现里和用例预期不一致，是实现侧的小 bug（当前没人用到根路径）。
- `module/send_email/parser_test.go` —— 硬编码了 `/home/jacko/...` 的绝对路径，
  在其它机器上必然找不到文件。

这两个和本次重构无关，先不处理。后续统一清理注释死代码那一轮一起修。
