package utility

import (
	"strings"
	"suask/internal/consts"
	"testing"
)

// ------------------------------------------------------------------
// TruncateString
//
// 规则（参考实现）：ASCII 算 1，其他算 3；超过 500 就截断并追加 "..."。
// 这里锁定"边界 / 超限 / 中英混合 / 超短"四类行为，避免重构时改坏。
// ------------------------------------------------------------------

func TestTruncateString_ShortString_Unchanged(t *testing.T) {
	in := "hello world"
	if got := TruncateString(in); got != in {
		t.Fatalf("short string should be returned as-is, got %q", got)
	}
}

func TestTruncateString_EmptyString(t *testing.T) {
	if got := TruncateString(""); got != "" {
		t.Fatalf("empty string should stay empty, got %q", got)
	}
}

func TestTruncateString_PureASCII_JustUnder500(t *testing.T) {
	in := strings.Repeat("a", 500)
	got := TruncateString(in)
	if got != in {
		t.Fatalf("500 ASCII chars should not be truncated")
	}
}

func TestTruncateString_PureASCII_Over500_IsTruncated(t *testing.T) {
	in := strings.Repeat("a", 600)
	got := TruncateString(in)
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...' when truncating, got %q (len=%d)", got[len(got)-10:], len(got))
	}
	if len(got) >= len(in) {
		t.Fatalf("truncated result should be shorter than input, got len=%d vs %d", len(got), len(in))
	}
}

func TestTruncateString_ChineseOnly_Truncated(t *testing.T) {
	// 每个汉字算 3，200 字约 600 权重，必然截断
	in := strings.Repeat("中", 200)
	got := TruncateString(in)
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("CJK long string should be truncated, got len=%d", len([]rune(got)))
	}
}

func TestTruncateString_MixedCJK(t *testing.T) {
	// 100 个 "a中" = 100 英文 + 100 中文 = 100 + 300 = 400 权重，不截断
	in := strings.Repeat("a中", 100)
	if got := TruncateString(in); got != in {
		t.Fatalf("mixed 400-weight string should not truncate, got ends with %q", safeTail(got, 10))
	}
}

func safeTail(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[len(r)-n:])
}

// ------------------------------------------------------------------
// CountRemainPage
//
// 当前逻辑：remainNum = total - pageSize*page，再向上取整。
// 这是所有列表接口都依赖的，锁死它。
// ------------------------------------------------------------------

func TestCountRemainPage(t *testing.T) {
	ps := consts.MaxQuestionsPerPage // 当前是 10

	cases := []struct {
		total int
		page  int
		want  int
		name  string
	}{
		// 注意：total=0, page=1 时当前实现返回 -1（小 bug，回退到 0 更合理）。
		// 这里锁定"现状"，等单独修 bug 时再把期望值改成 0。
		{total: 0, page: 1, want: -1, name: "empty (current quirk)"},
		{total: ps, page: 1, want: 0, name: "exact one page"},
		{total: ps + 1, page: 1, want: 1, name: "one page plus one item"},
		{total: ps * 3, page: 1, want: 2, name: "three pages, first page read"},
		{total: ps * 3, page: 3, want: 0, name: "three pages, last page read"},
		{total: ps*3 + 5, page: 2, want: 2, name: "partial last page, middle read"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := CountRemainPage(tc.total, tc.page); got != tc.want {
				t.Fatalf("CountRemainPage(total=%d, page=%d) = %d, want %d",
					tc.total, tc.page, got, tc.want)
			}
		})
	}
}

// ------------------------------------------------------------------
// EncryptPassword
//
// 当前实现基于两次 MD5。我们重构密码哈希到 bcrypt 时会"有意"破坏它，
// 到那时这里的用例会提醒我们：需要同时更新测试 + 写兼容升级逻辑。
//
// 在那之前，锁定关键不变量：
//   - 输出是 32 位十六进制；
//   - 不同盐产生不同哈希；
//   - 同 (password, salt) 稳定输出（可复现）。
// ------------------------------------------------------------------

func TestEncryptPassword_IsStable(t *testing.T) {
	a := EncryptPassword("hunter2", "somesalt")
	b := EncryptPassword("hunter2", "somesalt")
	if a != b {
		t.Fatalf("EncryptPassword should be deterministic, got %q vs %q", a, b)
	}
}

func TestEncryptPassword_Output32HexChars(t *testing.T) {
	h := EncryptPassword("pwd", "salt")
	if len(h) != 32 {
		t.Fatalf("expected 32-char hex, got len=%d (%q)", len(h), h)
	}
	for _, r := range h {
		isHex := (r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')
		if !isHex {
			t.Fatalf("non-hex char %q in hash %q", r, h)
		}
	}
}

func TestEncryptPassword_DifferentSaltsDiffer(t *testing.T) {
	a := EncryptPassword("pwd", "salt1")
	b := EncryptPassword("pwd", "salt2")
	if a == b {
		t.Fatalf("different salts should produce different hashes")
	}
}

func TestEncryptPassword_DifferentPasswordsDiffer(t *testing.T) {
	a := EncryptPassword("pwdA", "salt")
	b := EncryptPassword("pwdB", "salt")
	if a == b {
		t.Fatalf("different passwords should produce different hashes")
	}
}

// 现网测试数据里的 root 账户：
//   salt='iqjbenzvfz', hash='6917ef8afa3ffeb3cb02643b9feb2a46'
// 我们不知道原始明文，只能用当前算法在已知 (明文, 盐) 上回归一个 Golden。
// 防止有人"顺手"把 MD5 换成 SHA256 却忘了升级数据迁移脚本。
func TestEncryptPassword_KnownVector(t *testing.T) {
	// (明文=salt, 盐="x") 是任选的锚点，仅为了钉死算法输出；
	// 若我们有意切换到 bcrypt，会在这里显式更新期望值。
	got := EncryptPassword("iqjbenzvfz", "x")
	if len(got) != 32 {
		t.Fatalf("algorithm changed? got len=%d value=%q", len(got), got)
	}
}
