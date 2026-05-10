package login

import (
	"context"
	"testing"
)

// 原先这两个方法写的是 panic("implement me")，被误调会导致 goroutine 崩溃。
// 真正的业务逻辑在 controller 层，这里的 service 实现只要保证"安静且不崩"。
// 一旦有人以为这是空的就改成 panic，这两个用例会立刻告警。

func TestSLogin_Logout_DoesNotPanic(t *testing.T) {
	s := sLogin{}
	if err := s.Logout(context.Background()); err != nil {
		t.Fatalf("Logout should be a no-op, got err: %v", err)
	}
}

func TestSLogin_HeartBeats_DoesNotPanic(t *testing.T) {
	s := sLogin{}
	if err := s.HeartBeats(context.Background()); err != nil {
		t.Fatalf("HeartBeats should be a no-op, got err: %v", err)
	}
}
