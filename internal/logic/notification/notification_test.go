package notification

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"suask/internal/consts"
	"suask/internal/model"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func TestGetAndCountIgnoreDeletedNotifications(t *testing.T) {
	ctx := context.Background()
	setupNotificationTestDB(t, ctx)

	svc := New()
	out, err := svc.Get(ctx, model.GetNotificationsInput{UserId: 2})
	if err != nil {
		t.Fatalf("get notifications: %v", err)
	}
	if len(out.NewQuestion) != 1 {
		t.Fatalf("new question length = %d, want 1", len(out.NewQuestion))
	}
	if len(out.NewAnswer) != 1 {
		t.Fatalf("new answer length = %d, want 1", len(out.NewAnswer))
	}
	if len(out.NewReply) != 1 {
		t.Fatalf("new reply length = %d, want 1", len(out.NewReply))
	}
	question := out.NewQuestion[0]
	if question.UserId != 1 {
		t.Fatalf("new question user id = %d, want 1", question.UserId)
	}
	if question.UserName != "Student One" {
		t.Fatalf("new question user name = %q, want Student One", question.UserName)
	}
	if question.UserAvatar != consts.DefaultAvatarURL {
		t.Fatalf("new question user avatar = %q, want default avatar", question.UserAvatar)
	}

	count, err := svc.NewNotificationCount(ctx, model.NewNotificationCountInput{UserId: 2})
	if err != nil {
		t.Fatalf("count notifications: %v", err)
	}
	if count.NewQuestionCount != 1 {
		t.Fatalf("new question count = %d, want 1", count.NewQuestionCount)
	}
	if count.NewAnswerCount != 1 {
		t.Fatalf("new answer count = %d, want 1", count.NewAnswerCount)
	}
	if count.NewReplyCount != 1 {
		t.Fatalf("new reply count = %d, want 1", count.NewReplyCount)
	}
}

func TestGetNotificationFallbacksForDeletedContent(t *testing.T) {
	ctx := context.Background()
	setupNotificationTestDB(t, ctx)

	svc := New()
	out, err := svc.Get(ctx, model.GetNotificationsInput{UserId: 2})
	if err != nil {
		t.Fatalf("get notifications: %v", err)
	}

	answer := out.NewAnswer[0]
	if answer.AnswerContent != "" {
		t.Fatalf("deleted answer content = %q, want empty", answer.AnswerContent)
	}
	if answer.RespondentId != consts.DefaultUserId {
		t.Fatalf("deleted answer respondent id = %d, want %d", answer.RespondentId, consts.DefaultUserId)
	}
	if answer.RespondentName != consts.DefaultUserName {
		t.Fatalf("deleted answer respondent name = %q, want %q", answer.RespondentName, consts.DefaultUserName)
	}
	if answer.RespondentAvatar != consts.DefaultAvatarURL {
		t.Fatalf("deleted answer avatar = %q, want %q", answer.RespondentAvatar, consts.DefaultAvatarURL)
	}

	reply := out.NewReply[0]
	if reply.QuestionTitle != "" {
		t.Fatalf("deleted question title = %q, want empty", reply.QuestionTitle)
	}
	if reply.RespondentId != consts.DefaultUserId {
		t.Fatalf("teacher-owned reply respondent id = %d, want default user id", reply.RespondentId)
	}
	if reply.RespondentName != consts.DefaultUserName {
		t.Fatalf("teacher-owned reply respondent name = %q, want default user name", reply.RespondentName)
	}
	if reply.RespondentAvatar != consts.DefaultAvatarURL {
		t.Fatalf("teacher-owned reply avatar = %q, want default avatar", reply.RespondentAvatar)
	}
}

func TestNewNotificationCountOnlyUnreadActiveNotifications(t *testing.T) {
	ctx := context.Background()
	setupNotificationTestDB(t, ctx)

	if _, err := g.DB().Exec(
		ctx,
		`INSERT INTO notifications (id, user_id, question_id, answer_id, reply_to_id, type, is_read, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP)`,
		8,
		2,
		10,
		100,
		0,
		consts.NewAnswer,
		1,
	); err != nil {
		t.Fatalf("insert read notification: %v", err)
	}

	count, err := New().NewNotificationCount(ctx, model.NewNotificationCountInput{UserId: 2})
	if err != nil {
		t.Fatalf("count notifications: %v", err)
	}
	if count.NewAnswerCount != 1 {
		t.Fatalf("read notifications should not be counted, got new_answer_count=%d", count.NewAnswerCount)
	}
}

func TestUpdateAndDeleteNotificationRespectDeletedState(t *testing.T) {
	ctx := context.Background()
	setupNotificationTestDB(t, ctx)

	svc := New()
	updated, err := svc.Update(ctx, model.UpdateNotificationInput{Id: 1})
	if err != nil {
		t.Fatalf("update notification: %v", err)
	}
	if updated.Id != 1 || !updated.IsRead {
		t.Fatalf("update output = %+v, want id=1 is_read=true", updated)
	}
	if isRead := notificationIsRead(t, ctx, 1); !isRead {
		t.Fatalf("notification 1 should be marked read")
	}
	if _, err := svc.Update(ctx, model.UpdateNotificationInput{Id: 4}); err == nil {
		t.Fatalf("updating a deleted notification should fail")
	}

	if _, err := svc.Delete(ctx, model.DeleteNotificationInput{Id: 2}); err != nil {
		t.Fatalf("delete notification: %v", err)
	}
	if deleted := notificationIsDeleted(t, ctx, 2); !deleted {
		t.Fatalf("notification 2 should be soft deleted")
	}
	if _, err := svc.Delete(ctx, model.DeleteNotificationInput{Id: 2}); err == nil {
		t.Fatalf("deleting an already deleted notification should fail")
	}

	out, err := svc.Get(ctx, model.GetNotificationsInput{UserId: 2})
	if err != nil {
		t.Fatalf("get notifications after delete: %v", err)
	}
	if len(out.NewAnswer) != 0 {
		t.Fatalf("deleted new_answer notification should be hidden, got %d", len(out.NewAnswer))
	}
}

func TestUpdateAoQMarksOnlyActiveNotificationsRead(t *testing.T) {
	ctx := context.Background()
	setupNotificationTestDB(t, ctx)

	if _, err := New().UpdateAoQ(ctx, model.UpdateAoQInput{UserID: 2, QuestionID: 10}); err != nil {
		t.Fatalf("update notifications on question: %v", err)
	}

	if isRead := notificationIsRead(t, ctx, 1); !isRead {
		t.Fatalf("active question notification should be read")
	}
	if isRead := notificationIsRead(t, ctx, 2); !isRead {
		t.Fatalf("active answer notification should be read")
	}
	if isRead := notificationIsRead(t, ctx, 4); isRead {
		t.Fatalf("deleted question notification should not be changed")
	}
	if isRead := notificationIsRead(t, ctx, 5); isRead {
		t.Fatalf("deleted answer notification should not be changed")
	}
}

func setupNotificationTestDB(t *testing.T, ctx context.Context) {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), strings.ReplaceAll(t.Name(), "/", "_")+".db")
	if err := gdb.SetConfigGroup(gdb.DefaultGroupName, gdb.ConfigGroup{
		gdb.ConfigNode{
			Type: "sqlite",
			Link: fmt.Sprintf("sqlite::@file(%s)", dbPath),
		},
	}); err != nil {
		t.Fatalf("set test database config: %v", err)
	}

	db := g.DB()
	t.Cleanup(func() {
		if err := db.Close(ctx); err != nil {
			t.Fatalf("close test database: %v", err)
		}
	})

	statements := []string{
		`CREATE TABLE users (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT,
			role TEXT NOT NULL,
			nickname TEXT,
			avatar_file_id INTEGER,
			deleted_at DATETIME
		)`,
		`CREATE TABLE files (
			id INTEGER PRIMARY KEY,
			hash TEXT NOT NULL,
			name TEXT NOT NULL,
			uploader_id INTEGER,
			deleted_at DATETIME
		)`,
		`CREATE TABLE questions (
			id INTEGER PRIMARY KEY,
			src_user_id INTEGER NOT NULL,
			dst_user_id INTEGER NOT NULL,
			title TEXT NOT NULL,
			contents TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			views INTEGER NOT NULL DEFAULT 0,
			reply_cnt INTEGER NOT NULL DEFAULT 0,
			deleted_at DATETIME
		)`,
		`CREATE TABLE answers (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			question_id INTEGER,
			in_reply_to INTEGER DEFAULT 0,
			contents TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			upvotes INTEGER NOT NULL DEFAULT 0,
			deleted_at DATETIME
		)`,
		`CREATE TABLE notifications (
			id INTEGER PRIMARY KEY,
			user_id INTEGER NOT NULL,
			question_id INTEGER,
			reply_to_id INTEGER,
			answer_id INTEGER,
			type TEXT NOT NULL,
			is_read INTEGER NOT NULL DEFAULT 0,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			deleted_at DATETIME
		)`,
		`INSERT INTO files (id, hash, name, uploader_id) VALUES
			(20, '00112233445566778899aabbccddeeff', 'avatar.png', 3)`,
		`INSERT INTO users (id, name, email, role, nickname, avatar_file_id) VALUES
			(1, 'student1', 'student@example.com', 'student', 'Student One', 0),
			(2, 'teacher1', 'teacher@example.com', 'teacher', 'Teacher One', 0),
			(3, 'student2', 'student2@example.com', 'student', 'Student Two', 20)`,
		`INSERT INTO questions (id, src_user_id, dst_user_id, title, contents, views, reply_cnt, deleted_at) VALUES
			(10, 1, 2, 'Visible question', 'Visible body', 0, 1, NULL),
			(11, 1, 2, 'Deleted question', 'Deleted body', 0, 1, '2026-01-01 00:00:00')`,
		`INSERT INTO answers (id, user_id, question_id, contents, upvotes, deleted_at) VALUES
			(100, 3, 10, 'Deleted answer body', 0, '2026-01-01 00:00:00'),
			(101, 3, 11, 'Visible reply body', 0, NULL),
			(102, 1, 11, 'Original answer body', 0, NULL)`,
		`INSERT INTO notifications (id, user_id, question_id, answer_id, reply_to_id, type, is_read, created_at, deleted_at) VALUES
			(1, 2, 10, 0, 0, 'new_question', 0, '2026-01-01 00:00:03', NULL),
			(2, 2, 10, 100, 0, 'new_answer', 0, '2026-01-01 00:00:02', NULL),
			(3, 2, 11, 101, 102, 'new_reply', 0, '2026-01-01 00:00:01', NULL),
			(4, 2, 10, 0, 0, 'new_question', 0, '2026-01-01 00:00:04', '2026-01-01 00:01:00'),
			(5, 2, 10, 100, 0, 'new_answer', 0, '2026-01-01 00:00:05', '2026-01-01 00:01:00'),
			(6, 2, 11, 101, 102, 'new_reply', 0, '2026-01-01 00:00:06', '2026-01-01 00:01:00'),
			(7, 3, 10, 0, 0, 'new_question', 0, '2026-01-01 00:00:07', NULL)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(ctx, stmt); err != nil {
			t.Fatalf("exec test schema/fixture: %v\nSQL: %s", err, stmt)
		}
	}
}

func notificationIsRead(t *testing.T, ctx context.Context, id int) bool {
	t.Helper()

	var row struct {
		IsRead int `orm:"is_read"`
	}
	if err := g.DB().Ctx(ctx).
		Model("notifications").
		Unscoped().
		Fields("is_read").
		Where("id", id).
		Scan(&row); err != nil {
		t.Fatalf("query notification is_read: %v", err)
	}
	return row.IsRead == 1
}

func notificationIsDeleted(t *testing.T, ctx context.Context, id int) bool {
	t.Helper()

	var row struct {
		DeletedAt any `orm:"deleted_at"`
	}
	if err := g.DB().Ctx(ctx).
		Model("notifications").
		Unscoped().
		Fields("deleted_at").
		Where("id", id).
		Scan(&row); err != nil {
		t.Fatalf("query notification deleted_at: %v", err)
	}
	return row.DeletedAt != nil
}
