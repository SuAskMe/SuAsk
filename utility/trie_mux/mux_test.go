package triemux

import (
	"testing"
)

func TestNewTrieMux(t *testing.T) {
	trie := NewTrieMux()
	if trie == nil {
		t.Error("NewTrieMux() should not return nil")
	}
	if trie.root == nil {
		t.Error("NewTrieMux() should initialize root node")
	}
	if trie.root.children == nil {
		t.Error("NewTrieMux() should initialize root children map")
	}
}

func TestTrieMux_Insert(t *testing.T) {
	trie := NewTrieMux()

	// Test valid path insertion
	err := trie.Insert("/api/v1/users")
	if err != nil {
		t.Errorf("Insert(/api/v1/users) should not return error, got: %v", err)
	}

	// Test invalid path (not starting with /)
	err = trie.Insert("api/v1/users")
	if err == nil {
		t.Error("Insert(api/v1/users) should return error for path not starting with /")
	}

	// Test inserting empty segments
	err = trie.Insert("/")
	if err != nil {
		t.Errorf("Insert(/) should not return error, got: %v", err)
	}

	// Test multiple insertions
	paths := []string{"/api/v1/questions", "/api/v1/answers", "/api/v2/users"}
	for _, path := range paths {
		err = trie.Insert(path)
		if err != nil {
			t.Errorf("Insert(%s) should not return error, got: %v", path, err)
		}
	}
}

func TestTrieMux_HasPrefix(t *testing.T) {
	trie := NewTrieMux()

	// Insert some paths
	paths := []string{"/api/v1/users", "/api/v1/questions", "/api/v2/users", "/public/static"}
	for _, path := range paths {
		err := trie.Insert(path)
		if err != nil {
			t.Fatalf("Failed to insert path %s: %v", path, err)
		}
	}

	// Test existing paths
	existingPaths := []string{"/api/v1/users", "/api/v1/questions", "/api/v2/users", "/public/static"}
	for _, path := range existingPaths {
		if !trie.HasPrefix(path) {
			t.Errorf("HasPrefix(%s) should return true for existing path", path)
		}
	}

	// Test non-existing paths
	nonExistingPaths := []string{"/api/v1/answers", "/api/v3/users", "/private"}
	for _, path := range nonExistingPaths {
		if trie.HasPrefix(path) {
			t.Errorf("HasPrefix(%s) should return false for non-existing path", path)
		}
	}

	// Test invalid path (not starting with /)
	if trie.HasPrefix("api/v1/users") {
		t.Error("HasPrefix(api/v1/users) should return false for path not starting with /")
	}

	// Test root path
	if !trie.HasPrefix("/") {
		t.Error("HasPrefix(/) should return true as root path is always present")
	}
}

func TestTrieMux_GetSplitIndexFrom(t *testing.T) {
	trie := NewTrieMux()

	// Test various path splitting scenarios
	path := "/api/v1/users"
	index := trie.getSplitIndexFrom(&path, '/', 1)
	if index != 4 {
		t.Errorf("getSplitIndexFrom(/api/v1/users, '/', 1) should return 4, got %d", index)
	}

	path = "/api/v1/users"
	index = trie.getSplitIndexFrom(&path, '/', 5)
	if index != 7 {
		t.Errorf("getSplitIndexFrom(/api/v1/users, '/', 5) should return 7, got %d", index)
	}

	path = "/api/v1/users"
	index = trie.getSplitIndexFrom(&path, '/', 8)
	if index != 13 {
		t.Errorf("getSplitIndexFrom(/api/v1/users, '/', 8) should return 13, got %d", index)
	}

	// Test with character not in string
	path = "api-v1-users"
	index = trie.getSplitIndexFrom(&path, '/', 0)
	if index != 12 {
		t.Errorf("getSplitIndexFrom(api-v1-users, '/', 0) should return 12 (length of string), got %d", index)
	}
}
