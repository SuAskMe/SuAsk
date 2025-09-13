package send_email

import (
	"testing"
)

// 测试formatTmpl函数正常情况
func TestFormatTmpl(t *testing.T) {
	// 测试验证码模板
	// authCodeTmpl := `<html><body><p>您的验证码为：$code$</p></body></html>`
	authCodeTargets := []string{"$code$"}
	authCodeTmpl := loadTmplFile("/home/jacko/FILES/SuAsk/auth_code.tmpl", authCodeTargets)

	cache, err := formatTmpl(authCodeTmpl, authCodeTargets)
	if err != nil {
		t.Errorf("formatTmpl returned error: %v", err)
	}

	if cache == nil {
		t.Error("formatTmpl returned nil cache")
		return
	}

	if len(cache.pairs) != 1 {
		t.Errorf("Expected 1 pair, got %d", len(cache.pairs))
	}

	if len(cache.slice) != 2 {
		t.Errorf("Expected 2 slices, got %d", len(cache.slice))
	}

	// 测试消息模板

	messageTargets := []string{"$user$", "$type$", "$content$", "$url$"}
	messageTmpl := loadTmplFile("/home/jacko/FILES/SuAsk/msg.tmpl", messageTargets)
	cache, err = formatTmpl(messageTmpl, messageTargets)
	if err != nil {
		t.Errorf("formatTmpl returned error: %v", err)
	}

	if cache == nil {
		t.Error("formatTmpl returned nil cache")
		return
	}

	if len(cache.pairs) != 6 {
		t.Errorf("Expected 4 pairs, got %d", len(cache.pairs))
	}

	if len(cache.slice) != 7 {
		t.Errorf("Expected 5 slices, got %d", len(cache.slice))
	}

	// 测试重复占位符模板
	repeatTmpl := "$code$-$code$"
	cache, err = formatTmpl(repeatTmpl, authCodeTargets)
	if err != nil {
		t.Errorf("formatTmpl returned error: %v", err)
	}

	if cache == nil {
		t.Error("formatTmpl returned nil cache")
		return
	}

	if len(cache.pairs) != 2 {
		t.Errorf("Expected 2 pairs, got %d", len(cache.pairs))
	}

	if len(cache.slice) != 3 {
		t.Errorf("Expected 3 slices, got %d", len(cache.slice))
	}
}

// 测试formatTmpl函数错误情况
func TestFormatTmplError(t *testing.T) {
	tmpl := "Hello world"
	targets := []string{"$code$"}

	_, err := formatTmpl(tmpl, targets)
	if err == nil {
		t.Error("Expected error, but got nil")
	}
}

// 测试getAuthCode函数
func TestGetAuthCode(t *testing.T) {
	// 准备测试数据
	authCodeTmpl := `<html><body><p>您的验证码为：$code$</p></body></html>`
	authCodeTargets := []string{"$code$"}

	cache, err := formatTmpl(authCodeTmpl, authCodeTargets)
	if err != nil {
		t.Fatalf("formatTmpl failed: %v", err)
	}

	code := "123456"
	result := getAuthCode(cache, code)

	expected := `<html><body><p>您的验证码为：123456</p></body></html>`
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// 测试重复占位符
	repeatTmpl := "$code$-$code$"
	cache, err = formatTmpl(repeatTmpl, authCodeTargets)
	if err != nil {
		t.Fatalf("formatTmpl failed: %v", err)
	}

	result = getAuthCode(cache, "123")
	expected = "123-123"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

// 测试getMessage函数
func TestGetMessage(t *testing.T) {
	// 准备测试数据
	messageTmpl := `亲爱的$user$，您有新的消息！
[$type$] $content$
详情请访问以下链接：$url$
`
	messageTargets := []string{"$user$", "$type$", "$content$", "$url$"}

	cache, err := formatTmpl(messageTmpl, messageTargets)
	if err != nil {
		t.Fatalf("formatTmpl failed: %v", err)
	}

	notice := &Notice{
		User:    "张三",
		Type:    "提问",
		Content: "这是一个测试问题",
		URL:     "https://suask.example.com/question/1",
	}

	result := getMessage(cache, notice)

	expected := `亲爱的张三，您有新的消息！
[提问] 这是一个测试问题
详情请访问以下链接：https://suask.example.com/question/1
`
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
