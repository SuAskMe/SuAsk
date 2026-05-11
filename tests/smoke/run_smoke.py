"""
SuAsk backend 回归冒烟测试入口。

三种模式：
  seed       把测试账户写进 SQLite（幂等）
  snapshot   调用所有接口，录制 shape 到 tests/smoke/snapshots/
  verify     重新调用所有接口，和快照对比；有差异就非零退出

详细文档看同目录 README.md。
"""

from __future__ import annotations

import argparse
import json
import sys
import time
from pathlib import Path
from typing import Any, Callable

import requests

# 让脚本能直接 `python tests/smoke/run_smoke.py ...` 跑起来
_HERE = Path(__file__).resolve().parent
sys.path.insert(0, str(_HERE))

from seed import ensure_teacher, ensure_user  # noqa: E402
from shape import diff_shapes, shape_of  # noqa: E402


DEFAULT_BASE_URL = "http://127.0.0.1:8080"
DEFAULT_DB = Path("database/suask.db")
DEFAULT_USER = "smoke_student"
DEFAULT_EMAIL = "smoke_student@mail.sysu.edu.cn"
DEFAULT_PASSWORD = "Smoke-Pass-123!"

# 测试老师（smoke 里"发问"需要指定一个老师作为 dst_user_id；
# 这个账号也会被 seed 写进库，老师的提问箱权限置为 public 以方便跑 case）
DEFAULT_TEACHER = "smoke_teacher"
DEFAULT_TEACHER_EMAIL = "smoke_teacher@mail.sysu.edu.cn"
DEFAULT_TEACHER_PASSWORD = "Smoke-Teacher-123!"

SNAPSHOT_DIR = _HERE / "snapshots"


# ------------------------------------------------------------------
# HTTP 客户端
# ------------------------------------------------------------------


class Client:
    def __init__(self, base_url: str):
        self.base_url = base_url.rstrip("/")
        self.session = requests.Session()
        self.session.headers["Accept"] = "application/json"
        self._user_id: int | None = None

    def set_token(self, token: str) -> None:
        self.session.headers["Authorization"] = f"Bearer {token}"

    def set_user_id(self, uid: int) -> None:
        self._user_id = uid

    @property
    def user_id(self) -> int:
        if self._user_id is None:
            raise RuntimeError("user_id 还没设置；请先 login")
        return self._user_id

    def request(self, method: str, path: str, **kw) -> dict[str, Any]:
        url = f"{self.base_url}{path}"
        resp = self.session.request(method, url, timeout=10, **kw)
        # 后端始终返回 200 + JsonRes；业务错误放在 code 字段
        if not resp.headers.get("Content-Type", "").startswith("application/json"):
            raise AssertionError(
                f"{method} {path} 返回非 JSON：status={resp.status_code}, "
                f"body[:200]={resp.text[:200]!r}"
            )
        return resp.json()


# ------------------------------------------------------------------
# 登录：通过 /login 拿 token
# ------------------------------------------------------------------


def login(client: Client, name: str, password: str) -> int:
    body = client.request(
        "POST", "/login", json={"name": name, "password": password}
    )
    if body.get("code") != 0 or "data" not in body:
        raise AssertionError(f"登录失败: {body}")
    data = body["data"]
    token = data.get("token")
    uid = data.get("id")
    if not token or not uid:
        raise AssertionError(f"登录响应异常: {data}")
    client.set_token(token)
    client.set_user_id(int(uid))
    return int(uid)


# ------------------------------------------------------------------
# 每个被测接口封装成一个 case，返回的 JSON body 参与 shape 比对
# ------------------------------------------------------------------


Case = Callable[[Client], dict[str, Any]]


def case_login(client: Client) -> dict[str, Any]:
    # 专门再打一次 /login 只是为了快照结构；token 已经在 setup 阶段装好了
    return client.request(
        "POST",
        "/login",
        json={"name": DEFAULT_USER, "password": DEFAULT_PASSWORD},
    )


def case_user_self(client: Client) -> dict[str, Any]:
    return client.request("GET", "/user")


def case_user_by_id(client: Client) -> dict[str, Any]:
    return client.request(
        "GET", "/info/user", params={"id": client.user_id}
    )


def case_teacher_list(client: Client) -> dict[str, Any]:
    return client.request("GET", "/info/teacher")


def _get_teacher_id(client: Client) -> int:
    """从 /info/teacher 拿第一个老师的 id，供发问等 case 使用。"""
    body = client.request("GET", "/info/teacher")
    teachers = (body.get("data") or {}).get("teachers") or []
    if not teachers:
        raise AssertionError("没有老师数据，无法跑发问 case")
    return int(teachers[0]["id"])


def case_teacher_questions(client: Client) -> dict[str, Any]:
    """替代原来的 case_public_questions：现在所有问题都是问老师的。"""
    tid = _get_teacher_id(client)
    return client.request(
        "GET", "/questions/teacher",
        params={"sort_type": 0, "page": 1, "teacher_id": tid},
    )


def case_favorites_list(client: Client) -> dict[str, Any]:
    return client.request(
        "GET", "/favorites", params={"sort_type": 0, "page": 1}
    )


def case_history_list(client: Client) -> dict[str, Any]:
    return client.request(
        "GET", "/history", params={"sort_type": 0, "page": 1}
    )


def case_notification_list(client: Client) -> dict[str, Any]:
    return client.request(
        "GET", "/notification", params={"user_id": client.user_id}
    )


def case_notification_count(client: Client) -> dict[str, Any]:
    return client.request(
        "GET", "/notification/count", params={"user_id": client.user_id}
    )


def _pick_question_id(client: Client) -> int | None:
    """从老师问题列表里取第一条 id。"""
    tid = _get_teacher_id(client)
    body = client.request(
        "GET", "/questions/teacher",
        params={"sort_type": 0, "page": 1, "teacher_id": tid},
    )
    lst = (body.get("data") or {}).get("question_list") or []
    if not lst:
        return None
    return int(lst[0]["id"])


def case_question_detail(client: Client) -> dict[str, Any]:
    qid = _pick_question_id(client)
    if qid is None:
        # 列表空，造一条问老师的问题作为兜底
        tid = _get_teacher_id(client)
        r = client.request(
            "POST",
            "/questions/add",
            json={
                "dst_user_id": tid,
                "title": f"[smoke] detail-bootstrap {int(time.time())}",
                "content": "smoke bootstrap content",
                "is_private": False,
            },
        )
        qid = int(r["data"]["id"])
    return client.request("GET", "/answer", params={"question_id": qid})


def case_favorite_toggle(client: Client) -> dict[str, Any]:
    """收藏一次再取消，保持库干净。返回"第一次收藏"的响应做快照。"""
    qid = _pick_question_id(client)
    if qid is None:
        # 没有问题可收藏 —— 快照返回占位
        return {"code": 0, "message": "", "data": {"is_favorite": False}}
    first = client.request("POST", "/favorites", json={"question_id": qid})
    # 再调一次恢复成未收藏（如果第一次是 true，此次变 false）
    client.request("POST", "/favorites", json={"question_id": qid})
    return first


def case_add_question(client: Client) -> dict[str, Any]:
    """发一条问老师的问题（会留在库里，打上 [smoke] 标签方便识别）。"""
    tid = _get_teacher_id(client)
    title = f"[smoke] add-question {int(time.time())}"
    return client.request(
        "POST",
        "/questions/add",
        json={
            "dst_user_id": tid,
            "title": title,
            "content": "smoke regression test question, safe to delete",
            "is_private": False,
        },
    )


CASES: list[tuple[str, Case]] = [
    ("login", case_login),
    ("user_self", case_user_self),
    ("user_by_id", case_user_by_id),
    ("teacher_list", case_teacher_list),
    ("teacher_questions", case_teacher_questions),
    ("question_detail", case_question_detail),
    ("favorites_list", case_favorites_list),
    ("history_list", case_history_list),
    ("notification_list", case_notification_list),
    ("notification_count", case_notification_count),
    ("favorite_toggle", case_favorite_toggle),
    ("add_question", case_add_question),
]


# ------------------------------------------------------------------
# Snapshot / Verify
# ------------------------------------------------------------------


def snapshot_path(name: str) -> Path:
    return SNAPSHOT_DIR / f"{name}.json"


def run_snapshot(client: Client) -> None:
    SNAPSHOT_DIR.mkdir(parents=True, exist_ok=True)
    failures: list[str] = []
    for name, case in CASES:
        try:
            body = case(client)
        except Exception as exc:  # noqa: BLE001
            failures.append(f"[{name}] 调用失败: {exc}")
            continue
        shape = shape_of(body)
        snapshot_path(name).write_text(
            json.dumps(shape, ensure_ascii=False, indent=2, sort_keys=True) + "\n",
            encoding="utf-8",
        )
        print(f"[snapshot] {name}: OK")
    if failures:
        print("\n失败:", *failures, sep="\n  - ")
        sys.exit(1)


def run_verify(client: Client) -> None:
    failures: list[str] = []
    missing: list[str] = []
    for name, case in CASES:
        path = snapshot_path(name)
        if not path.exists():
            missing.append(name)
            continue
        try:
            body = case(client)
        except Exception as exc:  # noqa: BLE001
            failures.append(f"[{name}] 调用失败: {exc}")
            continue
        expected = json.loads(path.read_text(encoding="utf-8"))
        actual = shape_of(body)
        diffs = diff_shapes(expected, actual)
        if diffs:
            failures.append(f"[{name}] 结构变更:\n    " + "\n    ".join(diffs))
        else:
            print(f"[verify  ] {name}: OK")

    if missing:
        print(
            "\n以下用例没有快照（请先跑 `python tests/smoke/run_smoke.py snapshot`）:",
            *missing,
            sep="\n  - ",
        )
    if failures:
        print("\n发现回归:", *failures, sep="\n\n")
    if missing or failures:
        sys.exit(1)


# ------------------------------------------------------------------
# CLI
# ------------------------------------------------------------------


def build_parser() -> argparse.ArgumentParser:
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument("mode", choices=["seed", "snapshot", "verify"])
    p.add_argument("--base-url", default=DEFAULT_BASE_URL)
    p.add_argument("--db", default=str(DEFAULT_DB), type=Path)
    p.add_argument("--user", default=DEFAULT_USER)
    p.add_argument("--email", default=DEFAULT_EMAIL)
    p.add_argument("--password", default=DEFAULT_PASSWORD)
    return p


def main() -> None:
    args = build_parser().parse_args()

    # seed 只动数据库，不需要后端
    if args.mode == "seed":
        if not args.db.exists():
            print(f"[seed] 数据库文件不存在: {args.db}", file=sys.stderr)
            sys.exit(2)
        uid = ensure_user(args.db, args.user, args.email, args.password)
        tid = ensure_teacher(args.db, DEFAULT_TEACHER, DEFAULT_TEACHER_EMAIL, DEFAULT_TEACHER_PASSWORD)
        print(f"[seed] 测试学生就绪: id={uid}, name={args.user}")
        print(f"[seed] 测试老师就绪: id={tid}, name={DEFAULT_TEACHER}")
        return

    # snapshot / verify 都需要先登录；登录前同样保证账户存在
    if args.db.exists():
        ensure_user(args.db, args.user, args.email, args.password)
        ensure_teacher(args.db, DEFAULT_TEACHER, DEFAULT_TEACHER_EMAIL, DEFAULT_TEACHER_PASSWORD)
    else:
        print(
            f"[warn] 未找到数据库 {args.db}，假设账户已经由其他方式创建",
            file=sys.stderr,
        )

    client = Client(args.base_url)
    try:
        login(client, args.user, args.password)
    except requests.exceptions.ConnectionError:
        print(f"[error] 无法连接后端 {args.base_url}，请先启动服务", file=sys.stderr)
        sys.exit(2)

    if args.mode == "snapshot":
        run_snapshot(client)
    else:
        run_verify(client)


if __name__ == "__main__":
    main()
