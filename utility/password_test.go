package utility

import (
	"strings"
	"testing"
)

// HashPassword -----------------------------------------------------------------

func TestHashPassword_FormatLooksBcrypt(t *testing.T) {
	h, err := HashPassword("hunter2")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if len(h) < 59 || len(h) > 60 {
		t.Fatalf("bcrypt hash 应为 60 字节左右，got len=%d (%q)", len(h), h)
	}
	if !(strings.HasPrefix(h, "$2a$") || strings.HasPrefix(h, "$2b$") || strings.HasPrefix(h, "$2y$")) {
		t.Fatalf("bcrypt hash prefix 异常：%q", h)
	}
}

func TestHashPassword_NotDeterministic(t *testing.T) {
	// bcrypt 每次都会混入新盐，两次哈希应当不同
	a, _ := HashPassword("samepwd")
	b, _ := HashPassword("samepwd")
	if a == b {
		t.Fatalf("bcrypt 同明文两次应产生不同哈希，got 相同: %q", a)
	}
}

// IsLegacyHash -----------------------------------------------------------------

func TestIsLegacyHash(t *testing.T) {
	cases := []struct {
		in   string
		want bool
		note string
	}{
		{"", false, "empty"},
		{"6917ef8afa3ffeb3cb02643b9feb2a46", true, "真实老 MD5"},
		{"6917EF8AFA3FFEB3CB02643B9FEB2A46", true, "MD5 大写也算"},
		{"$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", false, "bcrypt 不是 legacy"},
		{"6917ef8afa3ffeb3cb02643b9feb2a4z", false, "长度 32 但含非 hex 字符"},
		{"abc", false, "长度不足"},
		{strings.Repeat("a", 64), false, "sha256 长度，不是 md5"},
	}
	for _, tc := range cases {
		t.Run(tc.note, func(t *testing.T) {
			if got := IsLegacyHash(tc.in); got != tc.want {
				t.Fatalf("IsLegacyHash(%q) = %v, want %v", tc.in, got, tc.want)
			}
		})
	}
}

// VerifyPassword ---------------------------------------------------------------

func TestVerifyPassword_Bcrypt_Match(t *testing.T) {
	hash, _ := HashPassword("secret")
	match, upgrade, err := VerifyPassword(hash, "", "secret")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !match {
		t.Fatal("正确密码应当 match=true")
	}
	if upgrade {
		t.Fatal("bcrypt 用户不应要求 upgrade")
	}
}

func TestVerifyPassword_Bcrypt_Mismatch(t *testing.T) {
	hash, _ := HashPassword("secret")
	match, _, err := VerifyPassword(hash, "", "wrong")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if match {
		t.Fatal("错误密码应当 match=false")
	}
}

func TestVerifyPassword_Legacy_MatchTriggersUpgrade(t *testing.T) {
	// 历史账户：手动算出 (password, salt) 对应的 md5(md5(pw)+md5(salt))，模拟 DB 里存的。
	const pw, salt = "hunter2", "somesalt"
	legacy := EncryptPassword(pw, salt)

	match, upgrade, err := VerifyPassword(legacy, salt, pw)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !match {
		t.Fatal("旧 MD5 正确密码应 match=true")
	}
	if !upgrade {
		t.Fatal("旧 MD5 正确密码必须返回 upgrade=true，否则不会升级到 bcrypt")
	}
}

func TestVerifyPassword_Legacy_Mismatch(t *testing.T) {
	const pw, salt = "hunter2", "somesalt"
	legacy := EncryptPassword(pw, salt)
	match, upgrade, err := VerifyPassword(legacy, salt, "wrong")
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if match {
		t.Fatal("错误密码应当 match=false")
	}
	if upgrade {
		t.Fatal("未匹配时不应要求 upgrade")
	}
}

func TestVerifyPassword_MalformedHash(t *testing.T) {
	// 既不像 md5（长度不对）也不是 bcrypt 的正确格式，应当返回 err
	_, _, err := VerifyPassword("this-is-not-a-valid-hash", "", "whatever")
	if err == nil {
		t.Fatal("格式异常的 hash 应当返回 err")
	}
}
