"""
Shape extraction + comparison.

"Shape" = 结构和类型，不记具体值。
这样生成的 snapshot 不会被时间戳、自增 ID、随机 salt 等噪声污染。

规则：
  - dict    → {key: shape(value)}；保留所有 key
  - list    → {"__list__": shape-of-first-element, "__empty__": bool}
              （我们只关心"列表里每项长什么样"，不关心长度）
  - str/int/float/bool/None → 类型名字符串 "string"/"int"/"float"/"bool"/"null"
"""

from __future__ import annotations

from typing import Any


def shape_of(value: Any) -> Any:
    if value is None:
        return "null"
    if isinstance(value, bool):  # 注意 bool 必须在 int 前判断
        return "bool"
    if isinstance(value, int):
        return "int"
    if isinstance(value, float):
        return "float"
    if isinstance(value, str):
        return "string"
    if isinstance(value, list):
        if not value:
            return {"__list__": "empty", "__empty__": True}
        # 取并集：列表里每个元素的 shape 合并，防止首元素不具代表性
        merged: Any = None
        for item in value:
            s = shape_of(item)
            merged = _merge_shapes(merged, s)
        return {"__list__": merged, "__empty__": False}
    if isinstance(value, dict):
        return {k: shape_of(v) for k, v in value.items()}
    return f"unknown:{type(value).__name__}"


_MISSING = object()  # 字典某侧缺少此 key 的哨兵，区别于 JSON null


def _merge_shapes(a: Any, b: Any) -> Any:
    """合并两份 shape：如果类型一致就按 dict 深合，否则标注成 union。

    约定：
      - Python None 表示"尚未初始化"（列表累积时的起点）
      - _MISSING 表示"这一侧字典根本没这个 key"
      - "null" 是 shape 层面的 JSON null
    """
    if a is None:
        return b
    if b is None:
        return a
    if a is _MISSING:
        return b
    if b is _MISSING:
        return a
    if a == b:
        return a
    if isinstance(a, dict) and isinstance(b, dict):
        keys = set(a) | set(b)
        return {
            k: _merge_shapes(a.get(k, _MISSING), b.get(k, _MISSING))
            for k in sorted(keys)
        }
    # 允许 null 和具体类型共存（典型场景：可空字段）
    if a == "null":
        return b
    if b == "null":
        return a
    return f"union({a}|{b})"


# ------------------------------------------------------------------
# diff
# ------------------------------------------------------------------


def diff_shapes(expected: Any, actual: Any, path: str = "$") -> list[str]:
    """返回 diff 消息列表；空列表代表一致。"""
    if expected == actual:
        return []
    if isinstance(expected, dict) and isinstance(actual, dict):
        msgs: list[str] = []
        exp_keys = set(expected)
        act_keys = set(actual)
        for k in sorted(exp_keys - act_keys):
            msgs.append(f"[-] {path}.{k}  (期望有该字段，响应里缺失)")
        for k in sorted(act_keys - exp_keys):
            msgs.append(f"[+] {path}.{k}  (响应里多出，快照未包含)")
        for k in sorted(exp_keys & act_keys):
            msgs.extend(diff_shapes(expected[k], actual[k], f"{path}.{k}"))
        return msgs
    return [f"[!] {path}: 快照={expected!r}, 当前={actual!r}"]
