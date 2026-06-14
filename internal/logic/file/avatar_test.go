package file

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"suask/internal/consts"
	"suask/internal/model"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
)

var (
	fileLogicTestDir  string
	fileLogicTestOnce sync.Once
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "suask-file-logic-*")
	if err != nil {
		panic(err)
	}
	fileLogicTestDir = dir
	code := m.Run()
	_ = g.DB().Close(context.Background())
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func TestFileGetMissingReturnsError(t *testing.T) {
	ctx := context.Background()
	setupFileLogicTest(t, ctx)

	if _, err := New().Get(ctx, model.FileGetInput{Id: 404}); err == nil {
		t.Fatalf("missing file should return error")
	}
}

func TestUploadAvatarCleansFileRecordWhenUserMissing(t *testing.T) {
	ctx := context.WithValue(context.Background(), consts.CtxId, 1)
	setupFileLogicTest(t, ctx)

	_, err := UploadAvatar(ctx, UploadAvatarInput{
		UserId: 999,
		File:   makeUploadFile(t, "avatar.png", pngBytes(64)),
	})
	if err == nil {
		t.Fatalf("upload avatar for missing user should fail")
	}
	if active := activeFileCount(t, ctx); active != 0 {
		t.Fatalf("active file count = %d, want 0", active)
	}
	if deleted := deletedFileCount(t, ctx); deleted != 1 {
		t.Fatalf("deleted file count = %d, want 1", deleted)
	}
}

func TestUploadAvatarUpdatesExistingUser(t *testing.T) {
	ctx := context.WithValue(context.Background(), consts.CtxId, 1)
	setupFileLogicTest(t, ctx)

	out, err := UploadAvatar(ctx, UploadAvatarInput{
		UserId: 1,
		File:   makeUploadFile(t, "avatar.png", pngBytes(64)),
	})
	if err != nil {
		t.Fatalf("upload avatar: %v", err)
	}
	if out.UserId != 1 {
		t.Fatalf("user id = %d, want 1", out.UserId)
	}
	if !strings.Contains(out.AvatarURL, "/upload/") {
		t.Fatalf("avatar url = %q, want upload url", out.AvatarURL)
	}
	if active := activeFileCount(t, ctx); active != 1 {
		t.Fatalf("active file count = %d, want 1", active)
	}
	if avatarFileId := userAvatarFileId(t, ctx, 1); avatarFileId == 0 {
		t.Fatalf("user avatar_file_id should be set")
	}
}

func setupFileLogicTest(t *testing.T, ctx context.Context) {
	t.Helper()

	fileLogicTestOnce.Do(func() {
		dbPath := filepath.Join(fileLogicTestDir, "file_logic.db")
		if err := gdb.SetConfigGroup(gdb.DefaultGroupName, gdb.ConfigGroup{
			gdb.ConfigNode{
				Type: "sqlite",
				Link: fmt.Sprintf("sqlite::@file(%s)", dbPath),
			},
		}); err != nil {
			t.Fatalf("set test database config: %v", err)
		}
		adapter, err := gcfg.NewAdapterContent(fmt.Sprintf(`
upload:
  path: "%s"
  prefix: "https://suask.test"
  max_minutes: 10
  max_bytes: 1048576
`, filepath.ToSlash(filepath.Join(fileLogicTestDir, "upload"))))
		if err != nil {
			t.Fatalf("create config adapter: %v", err)
		}
		g.Cfg().SetAdapter(adapter)
	})

	db := g.DB()
	statements := []string{
		`DROP TABLE IF EXISTS users`,
		`DROP TABLE IF EXISTS files`,
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT,
			role TEXT NOT NULL,
			nickname TEXT,
			avatar_file_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE files (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			hash BLOB NOT NULL,
			uploader_id INTEGER,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		)`,
		`INSERT INTO users (id, name, email, role, nickname)
		 VALUES (1, 'student1', 'student@example.com', 'student', 'Student One')`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(ctx, stmt); err != nil {
			t.Fatalf("exec test schema/fixture: %v\nSQL: %s", err, stmt)
		}
	}
}

func activeFileCount(t *testing.T, ctx context.Context) int {
	t.Helper()

	count, err := g.DB().Ctx(ctx).Model("files").Where("deleted_at IS NULL").Count()
	if err != nil {
		t.Fatalf("count active files: %v", err)
	}
	return count
}

func deletedFileCount(t *testing.T, ctx context.Context) int {
	t.Helper()

	count, err := g.DB().Ctx(ctx).Model("files").Unscoped().Where("deleted_at IS NOT NULL").Count()
	if err != nil {
		t.Fatalf("count deleted files: %v", err)
	}
	return count
}

func userAvatarFileId(t *testing.T, ctx context.Context, userId int) int {
	t.Helper()

	var row struct {
		AvatarFileId int `orm:"avatar_file_id"`
	}
	if err := g.DB().Ctx(ctx).
		Model("users").
		Fields("avatar_file_id").
		Where("id", userId).
		Scan(&row); err != nil {
		t.Fatalf("query user avatar file id: %v", err)
	}
	return row.AvatarFileId
}

func makeUploadFile(t *testing.T, filename string, payload []byte) *ghttp.UploadFile {
	t.Helper()

	var body bytes.Buffer
	w := multipart.NewWriter(&body)

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", `form-data; name="file"; filename="`+filename+`"`)
	header.Set("Content-Type", "application/octet-stream")
	part, err := w.CreatePart(header)
	if err != nil {
		t.Fatalf("create part: %v", err)
	}
	if _, err := part.Write(payload); err != nil {
		t.Fatalf("write part: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close writer: %v", err)
	}

	r := multipart.NewReader(&body, w.Boundary())
	form, err := r.ReadForm(int64(body.Len()) + 1024)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	files := form.File["file"]
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	return &ghttp.UploadFile{FileHeader: files[0]}
}

func pngBytes(extra int) []byte {
	head := []byte{0x89, 'P', 'N', 'G', 0x0D, 0x0A, 0x1A, 0x0A}
	return append(head, bytes.Repeat([]byte{0x00}, extra)...)
}
