package announcement

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"suask/internal/model"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

var (
	announcementTestDir  string
	announcementTestOnce sync.Once
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "suask-announcement-logic-*")
	if err != nil {
		panic(err)
	}
	announcementTestDir = dir
	code := m.Run()
	_ = g.DB().Close(context.Background())
	_ = os.RemoveAll(dir)
	os.Exit(code)
}

func TestAnnouncementUpdateUsesAffectedRowsAndCanClearExpiresAt(t *testing.T) {
	ctx := context.Background()
	setupAnnouncementTestDB(t, ctx)

	isPinned := false
	out, err := New().Update(ctx, model.AnnouncementUpdateInput{
		ID:             1,
		Title:          "Updated title",
		Content:        "Updated content",
		IsPinned:       &isPinned,
		ClearExpiresAt: true,
	})
	if err != nil {
		t.Fatalf("update announcement: %v", err)
	}
	if out.ID != 1 {
		t.Fatalf("update id = %d, want 1", out.ID)
	}

	row := announcementRow(t, ctx, 1)
	if row.Title != "Updated title" {
		t.Fatalf("title = %q, want Updated title", row.Title)
	}
	if row.Contents != "Updated content" {
		t.Fatalf("contents = %q, want Updated content", row.Contents)
	}
	if row.IsPinned != 0 {
		t.Fatalf("is_pinned = %d, want 0", row.IsPinned)
	}
	if row.ExpiresAt != nil {
		t.Fatalf("expires_at should be NULL after clear, got %v", row.ExpiresAt)
	}
	if row.DeletedAt != nil {
		t.Fatalf("update should not delete announcement")
	}

	if _, err := New().Update(ctx, model.AnnouncementUpdateInput{ID: 404, Title: "Missing"}); err == nil {
		t.Fatalf("updating missing announcement should fail")
	}
}

func TestAnnouncementUpdateWithNoFieldChangesOnlyChecksExistence(t *testing.T) {
	ctx := context.Background()
	setupAnnouncementTestDB(t, ctx)

	if _, err := New().Update(ctx, model.AnnouncementUpdateInput{ID: 1}); err != nil {
		t.Fatalf("empty update should succeed for existing announcement: %v", err)
	}
	if _, err := New().Update(ctx, model.AnnouncementUpdateInput{ID: 404}); err == nil {
		t.Fatalf("empty update should fail for missing announcement")
	}

	if err := New().Delete(ctx, model.AnnouncementDeleteInput{ID: 1}); err != nil {
		t.Fatalf("delete announcement: %v", err)
	}
	if _, err := New().Update(ctx, model.AnnouncementUpdateInput{ID: 1}); err == nil {
		t.Fatalf("empty update should fail for deleted announcement")
	}
}

func TestAnnouncementDeleteUsesSoftDeleteAndAffectedRows(t *testing.T) {
	ctx := context.Background()
	setupAnnouncementTestDB(t, ctx)

	if err := New().Delete(ctx, model.AnnouncementDeleteInput{ID: 1}); err != nil {
		t.Fatalf("delete announcement: %v", err)
	}
	if deletedAt := announcementRow(t, ctx, 1).DeletedAt; deletedAt == nil {
		t.Fatalf("deleted_at should be set")
	}
	if err := New().Delete(ctx, model.AnnouncementDeleteInput{ID: 1}); err == nil {
		t.Fatalf("deleting already deleted announcement should fail")
	}
	if err := New().Delete(ctx, model.AnnouncementDeleteInput{ID: 404}); err == nil {
		t.Fatalf("deleting missing announcement should fail")
	}
}

func TestAnnouncementCreateAndAddCommentUseStructuredData(t *testing.T) {
	ctx := context.Background()
	setupAnnouncementTestDB(t, ctx)

	expiresAt := gtime.NewFromStr("2026-02-01 00:00:00")
	created, err := New().Create(ctx, model.AnnouncementCreateInput{
		AuthorID:  1,
		Title:     "Created title",
		Content:   "Created content",
		IsPinned:  true,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		t.Fatalf("create announcement: %v", err)
	}
	row := announcementRow(t, ctx, created.ID)
	if row.AuthorId != 1 || row.Title != "Created title" || row.IsPinned != 1 {
		t.Fatalf("created announcement row = %+v", row)
	}

	replyTo := 10
	comment, err := New().AddComment(ctx, model.AnnouncementCommentInput{
		AnnouncementID: created.ID,
		UserID:         1,
		Content:        "Comment body",
		InReplyTo:      &replyTo,
	})
	if err != nil {
		t.Fatalf("add comment: %v", err)
	}
	if comment.ID == 0 {
		t.Fatalf("comment id should be set")
	}
	if got := announcementCommentCount(t, ctx, created.ID); got != 1 {
		t.Fatalf("comment count = %d, want 1", got)
	}
}

func setupAnnouncementTestDB(t *testing.T, ctx context.Context) {
	t.Helper()

	announcementTestOnce.Do(func() {
		dbPath := filepath.Join(announcementTestDir, "announcement_logic.db")
		if err := gdb.SetConfigGroup(gdb.DefaultGroupName, gdb.ConfigGroup{
			gdb.ConfigNode{
				Type: "sqlite",
				Link: fmt.Sprintf("sqlite::@file(%s)", dbPath),
			},
		}); err != nil {
			t.Fatalf("set test database config: %v", err)
		}
	})

	db := g.DB()
	statements := []string{
		`DROP TABLE IF EXISTS answers`,
		`DROP TABLE IF EXISTS announcements`,
		`DROP TABLE IF EXISTS users`,
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			nickname TEXT,
			deleted_at DATETIME
		)`,
		`CREATE TABLE announcements (
			id INTEGER PRIMARY KEY,
			author_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			contents TEXT NOT NULL,
			is_pinned INTEGER NOT NULL DEFAULT 0,
			published_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expires_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME,
			deleted_at DATETIME
		)`,
		`CREATE TABLE answers (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			question_id INTEGER,
			announcement_id INTEGER,
			in_reply_to INTEGER,
			contents TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			upvotes INTEGER NOT NULL DEFAULT 0,
			deleted_at DATETIME
		)`,
		`INSERT INTO users (id, nickname) VALUES (1, 'Admin One')`,
		`INSERT INTO announcements (id, author_id, title, contents, is_pinned, expires_at)
		 VALUES (1, 1, 'Original title', 'Original content', 1, '2026-01-01 00:00:00')`,
	}
	for _, stmt := range statements {
		if _, err := db.Exec(ctx, stmt); err != nil {
			t.Fatalf("exec test schema/fixture: %v\nSQL: %s", err, stmt)
		}
	}
}

type testAnnouncementRow struct {
	Id        int         `orm:"id"`
	AuthorId  int         `orm:"author_id"`
	Title     string      `orm:"title"`
	Contents  string      `orm:"contents"`
	IsPinned  int         `orm:"is_pinned"`
	ExpiresAt *gtime.Time `orm:"expires_at"`
	DeletedAt *gtime.Time `orm:"deleted_at"`
}

func announcementRow(t *testing.T, ctx context.Context, id int) testAnnouncementRow {
	t.Helper()

	var row testAnnouncementRow
	if err := g.DB().Ctx(ctx).
		Model("announcements").
		Unscoped().
		Where("id", id).
		Scan(&row); err != nil {
		t.Fatalf("query announcement row: %v", err)
	}
	return row
}

func announcementCommentCount(t *testing.T, ctx context.Context, announcementID int) int {
	t.Helper()

	count, err := g.DB().Ctx(ctx).
		Model("answers").
		Where("announcement_id", announcementID).
		Count()
	if err != nil {
		t.Fatalf("count announcement comments: %v", err)
	}
	return count
}
