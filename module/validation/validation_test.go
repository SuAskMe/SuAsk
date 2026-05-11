package validation

import (
	"suask/internal/model/entity"
	"sync"
	"testing"
)

// resetCache 给每个用例一个干净的 sync.Map，避免测试间互相污染。
func resetCache() {
	teacherCache = sync.Map{}
}

func TestUpdateTeacherPerm_PreservesPerm_AfterPriorLoad(t *testing.T) {
	resetCache()
	// 模拟：先前某处已经把完整的老师对象写进了缓存
	teacherCache.Store(11, &entity.Teachers{
		Id:   11,
		Perm: "public",
	})

	// 老师把提问箱权限改成了 private
	UpdateTeacherPerm(11, "苏玉鑫", "private")

	v, ok := teacherCache.Load(11)
	if !ok {
		t.Fatalf("teacher 11 should still be cached")
	}
	got, ok := v.(*entity.Teachers)
	if !ok {
		t.Fatalf("cached value must be *entity.Teachers, got %T", v)
	}
	if got.Perm != "private" {
		t.Fatalf("Perm 未被更新：got %q, want %q", got.Perm, "private")
	}
}

func TestUpdateTeacherPerm_CreatesEntry_WhenCacheEmpty(t *testing.T) {
	resetCache()

	UpdateTeacherPerm(42, "Alice", "protected")

	v, ok := teacherCache.Load(42)
	if !ok {
		t.Fatalf("teacher 42 应被写入缓存")
	}
	got := v.(*entity.Teachers)
	if got.Perm != "protected" {
		t.Fatalf("Perm 不正确: got %q", got.Perm)
	}
	if got.Id != 42 {
		t.Fatalf("Id 不正确: got %d", got.Id)
	}
}

func TestUpdateTeacherPerm_DoesNotMutateOldCachedPointer(t *testing.T) {
	resetCache()
	original := &entity.Teachers{Id: 7, Perm: "public"}
	teacherCache.Store(7, original)

	UpdateTeacherPerm(7, "Bob", "private")

	if original.Perm != "public" {
		t.Fatalf("旧缓存对象被意外就地修改：original.Perm=%q", original.Perm)
	}
	v, _ := teacherCache.Load(7)
	now := v.(*entity.Teachers)
	if now == original {
		t.Fatalf("缓存里仍然是旧指针，说明没做拷贝")
	}
	if now.Perm != "private" {
		t.Fatalf("新缓存对象 Perm 应为 private, got %q", now.Perm)
	}
}

func TestUpdateTeacherPerm_ReadPathReturnsFreshPerm(t *testing.T) {
	resetCache()
	teacherCache.Store(9, &entity.Teachers{Id: 9, Perm: "public"})

	UpdateTeacherPerm(9, "Carol", "private")

	perm, err := IsTeacher(nil, 9)
	if err != nil {
		t.Fatalf("IsTeacher err: %v", err)
	}
	if perm != "private" {
		t.Fatalf("IsTeacher 应看到 private, got %q", perm)
	}
}
