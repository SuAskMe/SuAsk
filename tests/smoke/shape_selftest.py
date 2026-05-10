"""
shape.py 自测：保证 shape 提取和 diff 没有低级 bug。
运行：python tests/smoke/shape_selftest.py
"""

from __future__ import annotations

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parent))

from shape import diff_shapes, shape_of


def _assert(cond: bool, msg: str) -> None:
    if not cond:
        print("FAIL:", msg, file=sys.stderr)
        sys.exit(1)


def test_scalars():
    _assert(shape_of("x") == "string", "string")
    _assert(shape_of(1) == "int", "int")
    _assert(shape_of(1.0) == "float", "float")
    _assert(shape_of(True) == "bool", "bool")
    _assert(shape_of(None) == "null", "null")


def test_dict():
    got = shape_of({"a": 1, "b": "hi", "c": None})
    want = {"a": "int", "b": "string", "c": "null"}
    _assert(got == want, f"dict: {got}")


def test_list_empty_vs_nonempty():
    _assert(
        shape_of([])["__empty__"] is True,
        "empty list must be marked empty",
    )
    s = shape_of([1, 2, 3])
    _assert(s["__list__"] == "int" and s["__empty__"] is False, f"int list: {s}")


def test_list_of_dicts_merge_keys():
    data = [{"a": 1}, {"b": 2}, {"a": 3, "b": 4}]
    s = shape_of(data)
    _assert(s["__list__"] == {"a": "int", "b": "int"}, f"merged dict list: {s}")


def test_diff_identical():
    _assert(diff_shapes({"a": "int"}, {"a": "int"}) == [], "identical")


def test_diff_missing_and_added():
    diffs = diff_shapes({"a": "int", "b": "string"}, {"a": "int", "c": "bool"})
    joined = "\n".join(diffs)
    _assert("$.b" in joined, "missing key should be flagged")
    _assert("$.c" in joined, "added key should be flagged")


def test_diff_type_change():
    diffs = diff_shapes({"a": "int"}, {"a": "string"})
    _assert(len(diffs) == 1 and "$.a" in diffs[0], f"type change: {diffs}")


def test_list_null_union():
    s = shape_of([{"x": 1}, {"x": None}])
    # null 和 int 应合并为 int
    _assert(s["__list__"] == {"x": "int"}, f"null+int should collapse to int: {s}")


def main():
    for fn in (
        test_scalars,
        test_dict,
        test_list_empty_vs_nonempty,
        test_list_of_dicts_merge_keys,
        test_diff_identical,
        test_diff_missing_and_added,
        test_diff_type_change,
        test_list_null_union,
    ):
        fn()
        print("OK  ", fn.__name__)
    print("\nall shape selftests passed")


if __name__ == "__main__":
    main()
