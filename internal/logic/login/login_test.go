package login

import (
	"context"
	"testing"
)

func TestSLogin_LogoutRequiresUser(t *testing.T) {
	s := sLogin{}
	if err := s.Logout(context.Background()); err == nil {
		t.Fatal("Logout should require a user in context")
	}
}

func TestSLogin_HeartBeatsRequiresUser(t *testing.T) {
	s := sLogin{}
	if err := s.HeartBeats(context.Background()); err == nil {
		t.Fatal("HeartBeats should require a user in context")
	}
}
