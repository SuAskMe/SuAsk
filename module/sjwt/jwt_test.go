package sjwt

import (
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// -----------------------------------------------------------------------------
// 测试辅助
// -----------------------------------------------------------------------------

// overrideKey 让测试用一个已知的 HMAC 密钥替换掉包级 jwtKey，
// 返回一个恢复函数；使用 defer 保证测试间互不污染。
func overrideKey(t *testing.T, key string) func() {
	t.Helper()
	old := jwtKey
	jwtKey = key
	return func() { jwtKey = old }
}

// signedToken 用给定 method + key 签一个 JwtClaims，
// 方便造"签名算法不对"或"密钥不对"的场景。
func signedToken(t *testing.T, method jwt.SigningMethod, key any, claims *JwtClaims) string {
	t.Helper()
	tok := jwt.NewWithClaims(method, claims)
	out, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign token err: %v", err)
	}
	return out
}

func baseClaims(userID int, exp time.Duration) *JwtClaims {
	return &JwtClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "SuAsk",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
		UserID: userID,
	}
}

// -----------------------------------------------------------------------------
// 正向
// -----------------------------------------------------------------------------

func TestParseToken_ValidToken(t *testing.T) {
	defer overrideKey(t, "unit-test-secret")()

	token := signedToken(t,
		jwt.SigningMethodHS256,
		[]byte(jwtKey),
		baseClaims(42, time.Hour),
	)

	claims, err := ParseToken(token)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if claims == nil {
		t.Fatal("claims should not be nil for valid token")
	}
	if claims.UserID != 42 {
		t.Fatalf("UserID mismatch: got %d, want 42", claims.UserID)
	}
}

// -----------------------------------------------------------------------------
// 反向：这些场景之前可能返回 (nil, nil)，上层会误判"成功"。
// -----------------------------------------------------------------------------

func TestParseToken_TamperedToken(t *testing.T) {
	defer overrideKey(t, "unit-test-secret")()

	token := signedToken(t,
		jwt.SigningMethodHS256,
		[]byte(jwtKey),
		baseClaims(42, time.Hour),
	)
	// 篡改最后一个字符
	tampered := token[:len(token)-1] + "A"

	claims, err := ParseToken(tampered)
	if err == nil {
		t.Fatal("tampered token must return error")
	}
	if claims != nil {
		t.Fatalf("tampered token must return nil claims, got %+v", claims)
	}
}

func TestParseToken_ExpiredToken(t *testing.T) {
	defer overrideKey(t, "unit-test-secret")()

	// 1 小时前就过期
	token := signedToken(t,
		jwt.SigningMethodHS256,
		[]byte(jwtKey),
		baseClaims(42, -time.Hour),
	)
	claims, err := ParseToken(token)
	if err == nil {
		t.Fatal("expired token must return error")
	}
	if claims != nil {
		t.Fatalf("expired token must return nil claims, got %+v", claims)
	}
}

func TestParseToken_WrongKey(t *testing.T) {
	// 用 key A 签发，用 key B 解析
	restore := overrideKey(t, "key-A")
	token := signedToken(t,
		jwt.SigningMethodHS256,
		[]byte(jwtKey),
		baseClaims(42, time.Hour),
	)
	restore()

	defer overrideKey(t, "key-B-different")()
	claims, err := ParseToken(token)
	if err == nil {
		t.Fatal("token signed with different key must fail")
	}
	if claims != nil {
		t.Fatalf("wrong-key token must return nil claims, got %+v", claims)
	}
}

// 防御 alg confusion：即使攻击者改写 header 声明 alg=none 且不带签名，
// 也绝不能被认作合法。
func TestParseToken_RejectsNoneAlg(t *testing.T) {
	defer overrideKey(t, "unit-test-secret")()

	noneTok := jwt.NewWithClaims(jwt.SigningMethodNone, baseClaims(42, time.Hour))
	out, err := noneTok.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		t.Fatalf("sign none err: %v", err)
	}

	claims, err := ParseToken(out)
	if err == nil {
		t.Fatal("alg=none token must be rejected")
	}
	if claims != nil {
		t.Fatalf("alg=none must return nil claims, got %+v", claims)
	}
	if !strings.Contains(err.Error(), "unexpected signing method") {
		t.Logf("(提示) 错误信息未包含 'unexpected signing method'，当前: %v", err)
	}
}

func TestParseToken_Garbage(t *testing.T) {
	defer overrideKey(t, "unit-test-secret")()

	cases := []string{
		"",
		"not-a-jwt",
		"a.b",         // 段数不对
		"a.b.c.d.e",   // 段数不对
		"....",        // 全是分隔符
		"header.payload.signature",
	}
	for _, in := range cases {
		claims, err := ParseToken(in)
		if err == nil {
			t.Errorf("garbage input %q should error", in)
		}
		if claims != nil {
			t.Errorf("garbage input %q should return nil claims, got %+v", in, claims)
		}
	}
}

// -----------------------------------------------------------------------------
// GenerateToken 与 ParseToken 的闭环
// -----------------------------------------------------------------------------

func TestGenerate_Parse_RoundTrip(t *testing.T) {
	defer overrideKey(t, "unit-test-secret")()
	// 让 jwtExpire 不为 0，避免生成的 token 瞬间过期
	oldExpire := jwtExpire
	jwtExpire = 1
	defer func() { jwtExpire = oldExpire }()

	tok, err := GenerateToken(1001)
	if err != nil {
		t.Fatalf("generate err: %v", err)
	}
	claims, err := ParseToken(tok)
	if err != nil {
		t.Fatalf("parse err: %v", err)
	}
	if claims.UserID != 1001 {
		t.Fatalf("UserID round-trip mismatch: got %d", claims.UserID)
	}
	if claims.Issuer != "SuAsk" {
		t.Fatalf("Issuer mismatch: got %q", claims.Issuer)
	}
}
