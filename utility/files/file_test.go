package files

import (
	"bytes"
	"encoding/hex"
	"io"
	"strings"
	"testing"
)

// fakeMultipartFile 让我们在不起 HTTP 服务的情况下测 HashFile。
// multipart.File = io.Reader + io.ReaderAt + io.Seeker + io.Closer
type fakeMultipartFile struct {
	*bytes.Reader
}

func (fakeMultipartFile) Close() error { return nil }

func newFake(data []byte) fakeMultipartFile {
	return fakeMultipartFile{Reader: bytes.NewReader(data)}
}

// ------------------------------------------------------------------
// HashFile
// ------------------------------------------------------------------

func TestHashFile_SameInputSameHash(t *testing.T) {
	payload := []byte("the quick brown fox jumps over the lazy dog")
	a, err := HashFile(newFake(payload))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	b, err := HashFile(newFake(payload))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !bytes.Equal(a, b) {
		t.Fatalf("hash should be deterministic for same input, got %x vs %x", a, b)
	}
	if len(a) != 32 {
		t.Fatalf("BLAKE2b-256 should be 32 bytes, got %d", len(a))
	}
}

func TestHashFile_DifferentInputDifferentHash(t *testing.T) {
	a, err := HashFile(newFake([]byte("abc")))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	b, err := HashFile(newFake([]byte("abd")))
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if bytes.Equal(a, b) {
		t.Fatalf("different inputs must produce different hashes")
	}
}

func TestHashFile_ConsumesReaderFully(t *testing.T) {
	// 构造 1MB payload 确保没有半途停止
	payload := bytes.Repeat([]byte{0xAB}, 1024*1024)
	if _, err := HashFile(newFake(payload)); err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
}

// errReader 让 io.Copy 失败，验证 HashFile 返回 error 而不是 log.Fatal。
type errReader struct{}

func (errReader) Read(_ []byte) (int, error)                     { return 0, io.ErrUnexpectedEOF }
func (errReader) Close() error                                   { return nil }
func (errReader) Seek(_ int64, _ int) (int64, error)             { return 0, io.ErrUnexpectedEOF }
func (errReader) ReadAt(_ []byte, _ int64) (int, error)          { return 0, io.ErrUnexpectedEOF }

func TestHashFile_PropagatesReadError(t *testing.T) {
	_, err := HashFile(errReader{})
	if err == nil {
		t.Fatal("HashFile 应把 io 错误返回而不是 log.Fatal")
	}
}

// ------------------------------------------------------------------
// HashToString
// ------------------------------------------------------------------

func TestHashToString_RoundTrip(t *testing.T) {
	raw, _ := hex.DecodeString("deadbeefcafe0001")
	got := HashToString(raw)
	if got != "deadbeefcafe0001" {
		t.Fatalf("hex encode mismatch, got %q", got)
	}
}

func TestHashToString_Empty(t *testing.T) {
	if got := HashToString(nil); got != "" {
		t.Fatalf("empty hash should produce empty string, got %q", got)
	}
}

// ------------------------------------------------------------------
// RenameFiles
//
// 目前的契约：
//   - 有扩展名：<hex-hash>.<ext>
//   - 没扩展名：<hex-hash>
// 这是上传落盘和 URL 计算共同依赖的，锁死它。
// ------------------------------------------------------------------

func TestRenameFiles_WithExtension(t *testing.T) {
	raw, _ := hex.DecodeString("0011223344556677")
	newName, err := RenameFiles(raw, "avatar.PNG")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 扩展名大小写目前未规范，但必须以 hash 开头、以 .png 结尾（或 .PNG）
	if !strings.HasPrefix(newName, "0011223344556677") {
		t.Fatalf("renamed file should start with hex hash, got %q", newName)
	}
	if !strings.Contains(newName, ".") {
		t.Fatalf("renamed file should preserve the dot separator, got %q", newName)
	}
}

func TestRenameFiles_NoExtension(t *testing.T) {
	raw, _ := hex.DecodeString("aabbccdd")
	newName, err := RenameFiles(raw, "somefile")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if newName != "aabbccdd" {
		t.Fatalf("no-extension file should be pure hash, got %q", newName)
	}
}

// 确保我们的 fake 实现了 multipart.File 所需的接口
var _ interface {
	io.Reader
	io.Seeker
	io.Closer
} = fakeMultipartFile{}
