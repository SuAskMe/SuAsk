package admin

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	v1 "suask/api/admin/v1"

	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

func TestResolveAdminDeletedStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		deletedStatus  string
		includeDeleted bool
		want           string
		wantErr        bool
	}{
		{name: "default excludes deleted", want: adminDeletedStatusUndeleted},
		{name: "compat include deleted", includeDeleted: true, want: adminDeletedStatusAll},
		{name: "explicit all", deletedStatus: adminDeletedStatusAll, want: adminDeletedStatusAll},
		{name: "explicit deleted", deletedStatus: adminDeletedStatusDeleted, want: adminDeletedStatusDeleted},
		{name: "explicit undeleted", deletedStatus: adminDeletedStatusUndeleted, want: adminDeletedStatusUndeleted},
		{name: "invalid", deletedStatus: "visible", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := resolveAdminDeletedStatus(tt.deletedStatus, tt.includeDeleted)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("status = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestQuestionAnswerDeleteRestoreMaintainsCountsAndFilters(t *testing.T) {
	ctx := context.Background()
	setupAdminLogicTestDB(t, ctx)

	detail, err := GetQuestionDetail(ctx, &v1.GetQuestionDetailReq{
		Id:            10,
		DeletedStatus: adminDeletedStatusUndeleted,
	})
	if err != nil {
		t.Fatalf("get initial detail: %v", err)
	}
	assertAdminQuestionCounts(t, detail, 1, 1, 1)
	if detail.Answers[0].IsDeleted {
		t.Fatalf("initial answer should be visible")
	}

	if _, err = DeleteQuestionAnswer(ctx, 10, 100); err != nil {
		t.Fatalf("delete answer: %v", err)
	}
	detail, err = GetQuestionDetail(ctx, &v1.GetQuestionDetailReq{
		Id:            10,
		DeletedStatus: adminDeletedStatusUndeleted,
	})
	if err != nil {
		t.Fatalf("get detail after delete: %v", err)
	}
	assertAdminQuestionCounts(t, detail, 0, 0, 0)

	detail, err = GetQuestionDetail(ctx, &v1.GetQuestionDetailReq{
		Id:             10,
		IncludeDeleted: true,
	})
	if err != nil {
		t.Fatalf("get detail with deleted answer: %v", err)
	}
	assertAdminQuestionCounts(t, detail, 0, 0, 1)
	if !detail.Answers[0].IsDeleted {
		t.Fatalf("deleted answer should be returned as deleted when include_deleted is true")
	}

	if _, err = RestoreQuestionAnswer(ctx, 10, 100); err != nil {
		t.Fatalf("restore answer: %v", err)
	}
	detail, err = GetQuestionDetail(ctx, &v1.GetQuestionDetailReq{
		Id:            10,
		DeletedStatus: adminDeletedStatusUndeleted,
	})
	if err != nil {
		t.Fatalf("get detail after restore: %v", err)
	}
	assertAdminQuestionCounts(t, detail, 1, 1, 1)
	if detail.Answers[0].IsDeleted {
		t.Fatalf("restored answer should be visible")
	}
}

func TestQuestionAnswerRepeatedDeleteRestoreDoesNotDriftCounts(t *testing.T) {
	ctx := context.Background()
	setupAdminLogicTestDB(t, ctx)

	if _, err := DeleteQuestionAnswer(ctx, 10, 100); err != nil {
		t.Fatalf("delete answer: %v", err)
	}
	if _, err := DeleteQuestionAnswer(ctx, 10, 100); err == nil {
		t.Fatalf("repeated delete should fail")
	}
	detail, err := GetQuestionDetail(ctx, &v1.GetQuestionDetailReq{
		Id:             10,
		IncludeDeleted: true,
	})
	if err != nil {
		t.Fatalf("get detail after repeated delete: %v", err)
	}
	assertAdminQuestionCounts(t, detail, 0, 0, 1)

	if _, err := RestoreQuestionAnswer(ctx, 10, 100); err != nil {
		t.Fatalf("restore answer: %v", err)
	}
	if _, err := RestoreQuestionAnswer(ctx, 10, 100); err == nil {
		t.Fatalf("repeated restore should fail")
	}
	detail, err = GetQuestionDetail(ctx, &v1.GetQuestionDetailReq{
		Id:            10,
		DeletedStatus: adminDeletedStatusUndeleted,
	})
	if err != nil {
		t.Fatalf("get detail after repeated restore: %v", err)
	}
	assertAdminQuestionCounts(t, detail, 1, 1, 1)
}

func TestListQuestionsStatusFilters(t *testing.T) {
	ctx := context.Background()
	setupAdminLogicTestDB(t, ctx)

	insertQuestion(t, ctx, 11, "Visible unanswered", 0, "")
	insertQuestion(t, ctx, 12, "Deleted question", 1, "2026-01-01 00:00:00")
	insertAnswer(t, ctx, 101, 2, 12, "Deleted question answer")

	tests := []struct {
		name       string
		status     string
		wantTotal  int
		wantIDs    map[int]bool
		wantErr    bool
		wantRemain int
	}{
		{
			name:      "blank defaults to undeleted questions",
			wantTotal: 2,
			wantIDs:   map[int]bool{10: true, 11: true},
		},
		{
			name:      "all means all undeleted questions",
			status:    adminQuestionListStatusAll,
			wantTotal: 2,
			wantIDs:   map[int]bool{10: true, 11: true},
		},
		{
			name:      "answered only counts visible answers",
			status:    adminQuestionListStatusAnswered,
			wantTotal: 1,
			wantIDs:   map[int]bool{10: true},
		},
		{
			name:      "unanswered excludes deleted questions",
			status:    adminQuestionListStatusUnanswered,
			wantTotal: 1,
			wantIDs:   map[int]bool{11: true},
		},
		{
			name:      "deleted returns only deleted questions",
			status:    adminQuestionListStatusDeleted,
			wantTotal: 1,
			wantIDs:   map[int]bool{12: true},
		},
		{
			name:    "invalid status fails",
			status:  "visible",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := ListQuestions(ctx, &v1.ListQuestionsReq{
				Page:   1,
				Status: tt.status,
			})
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("list questions: %v", err)
			}
			if res.Total != tt.wantTotal {
				t.Fatalf("total = %d, want %d", res.Total, tt.wantTotal)
			}
			if res.RemainPage != tt.wantRemain {
				t.Fatalf("remain_page = %d, want %d", res.RemainPage, tt.wantRemain)
			}
			if len(res.List) != tt.wantTotal {
				t.Fatalf("list length = %d, want %d", len(res.List), tt.wantTotal)
			}
			for _, item := range res.List {
				if !tt.wantIDs[item.Id] {
					t.Fatalf("unexpected question id %d in status %q", item.Id, tt.status)
				}
			}
		})
	}
}

func TestQuestionDetailDeletedStatusFilters(t *testing.T) {
	ctx := context.Background()
	setupAdminLogicTestDB(t, ctx)

	if _, err := DeleteQuestionAnswer(ctx, 10, 100); err != nil {
		t.Fatalf("delete answer: %v", err)
	}

	tests := []struct {
		name          string
		deletedStatus string
		wantLen       int
		wantDeleted   bool
	}{
		{name: "undeleted hides deleted answers", deletedStatus: adminDeletedStatusUndeleted, wantLen: 0},
		{name: "deleted returns only deleted answers", deletedStatus: adminDeletedStatusDeleted, wantLen: 1, wantDeleted: true},
		{name: "all returns deleted answers", deletedStatus: adminDeletedStatusAll, wantLen: 1, wantDeleted: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			detail, err := GetQuestionDetail(ctx, &v1.GetQuestionDetailReq{
				Id:            10,
				DeletedStatus: tt.deletedStatus,
			})
			if err != nil {
				t.Fatalf("get detail: %v", err)
			}
			if len(detail.Answers) != tt.wantLen {
				t.Fatalf("answer length = %d, want %d", len(detail.Answers), tt.wantLen)
			}
			if tt.wantLen > 0 && detail.Answers[0].IsDeleted != tt.wantDeleted {
				t.Fatalf("answer deleted = %v, want %v", detail.Answers[0].IsDeleted, tt.wantDeleted)
			}
		})
	}
}

func TestRestoreQuestionClearsDeletedState(t *testing.T) {
	ctx := context.Background()
	setupAdminLogicTestDB(t, ctx)
	insertQuestion(t, ctx, 12, "Deleted question", 0, "2026-01-01 00:00:00")

	if _, err := RestoreQuestion(ctx, 12); err != nil {
		t.Fatalf("restore question: %v", err)
	}

	res, err := ListQuestions(ctx, &v1.ListQuestionsReq{
		Page:   1,
		Status: adminQuestionListStatusAll,
	})
	if err != nil {
		t.Fatalf("list questions: %v", err)
	}
	found := false
	for _, item := range res.List {
		if item.Id == 12 {
			found = true
			if item.IsDeleted {
				t.Fatalf("restored question should not be marked deleted")
			}
		}
	}
	if !found {
		t.Fatalf("restored question should appear in undeleted list")
	}

	if _, err = RestoreQuestion(ctx, 12); err == nil {
		t.Fatalf("restoring an undeleted question should fail")
	}
}

func setupAdminLogicTestDB(t *testing.T, ctx context.Context) {
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
		`INSERT INTO users (id, name, email, role, nickname) VALUES
			(1, 'student1', 'student@example.com', 'student', 'Student One'),
			(2, 'teacher1', 'teacher@example.com', 'teacher', 'Teacher One')`,
		`INSERT INTO questions (id, src_user_id, dst_user_id, title, contents, views, reply_cnt)
			VALUES (10, 1, 2, 'Question title', 'Question body', 7, 1)`,
		`INSERT INTO answers (id, user_id, question_id, contents, upvotes)
			VALUES (100, 2, 10, 'Answer body', 3)`,
	}

	for _, stmt := range statements {
		if _, err := db.Exec(ctx, stmt); err != nil {
			t.Fatalf("exec test schema/fixture: %v\nSQL: %s", err, stmt)
		}
	}
}

func insertQuestion(t *testing.T, ctx context.Context, id int, title string, replyCnt int, deletedAt string) {
	t.Helper()

	var deletedValue any
	if deletedAt != "" {
		deletedValue = deletedAt
	}
	if _, err := g.DB().Exec(
		ctx,
		`INSERT INTO questions (id, src_user_id, dst_user_id, title, contents, views, reply_cnt, deleted_at)
		 VALUES (?, 1, 2, ?, 'Question body', 0, ?, ?)`,
		id,
		title,
		replyCnt,
		deletedValue,
	); err != nil {
		t.Fatalf("insert question fixture: %v", err)
	}
}

func insertAnswer(t *testing.T, ctx context.Context, id int, userId int, questionId int, contents string) {
	t.Helper()

	if _, err := g.DB().Exec(
		ctx,
		`INSERT INTO answers (id, user_id, question_id, contents, upvotes)
		 VALUES (?, ?, ?, ?, 0)`,
		id,
		userId,
		questionId,
		contents,
	); err != nil {
		t.Fatalf("insert answer fixture: %v", err)
	}
}

func assertAdminQuestionCounts(t *testing.T, detail *v1.GetQuestionDetailRes, replyCnt int, answerCount int, answerLen int) {
	t.Helper()

	if detail.Question.ReplyCnt != replyCnt {
		t.Fatalf("reply_cnt = %d, want %d", detail.Question.ReplyCnt, replyCnt)
	}
	if detail.Question.AnswerCount != answerCount {
		t.Fatalf("answer_count = %d, want %d", detail.Question.AnswerCount, answerCount)
	}
	if len(detail.Answers) != answerLen {
		t.Fatalf("answer length = %d, want %d", len(detail.Answers), answerLen)
	}
}
