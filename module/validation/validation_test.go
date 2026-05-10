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

func TestUpdateTeacherPerm_PreservesName_AfterPriorLoad(t *testing.T) {
	resetCache()
	// 模拟：先前某处已经把完整的老师对象写进了缓存
	teacherCache.Store(11, &entity.Teachers{
		Id:    11,
		Name:  "苏玉鑫",
		Email: "suyx35@mail.sysu.edu.cn",
		Perm:  "public",
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
	// 这里是本次修复的核心断言：Name / Email 不能被清掉
	if got.Name != "苏玉鑫" {
		t.Fatalf("Name 被清掉了：got %q, want %q", got.Name, "苏玉鑫")
	}
	if got.Email != "suyx35@mail.sysu.edu.cn" {
		t.Fatalf("Email 被清掉了：got %q", got.Email)
	}
}

func TestUpdateTeacherPerm_CreatesEntry_WhenCacheEmpty(t *testing.T) {
	resetCache()

	// 缓存里没有 42 号老师，这时调 UpdatePerm 也应该能建立条目
	UpdateTeacherPerm(42, "Alice", "protected")

	v, ok := teacherCache.Load(42)
	if !ok {
		t.Fatalf("teacher 42 应被写入缓存")
	}
	got := v.(*entity.Teachers)
	if got.Perm != "protected" {
		t.Fatalf("Perm 不正确: got %q", got.Perm)
	}
	if got.Name != "Alice" {
		t.Fatalf("Name 不正确: got %q", got.Name)
	}
	if got.Id != 42 {
		t.Fatalf("Id 不正确: got %d", got.Id)
	}
}

func TestUpdateTeacherPerm_DoesNotMutateOldCachedPointer(t *testing.T) {
	// 回归：其他 goroutine 在 UpdatePerm 之前拿到的指针不应被别名修改。
	resetCache()
	original := &entity.Teachers{Id: 7, Name: "Bob", Perm: "public"}
	teacherCache.Store(7, original)

	UpdateTeacherPerm(7, "Bob", "private")

	if original.Perm != "public" {
		t.Fatalf("旧缓存对象被意外就地修改：original.Perm=%q", original.Perm)
	}
	// 新读出来的应是更新过的
	v, _ := teacherCache.Load(7)
	now := v.(*entity.Teachers)
	if now == original {
		t.Fatalf("缓存里仍然是旧指针，说明没做拷贝")
	}
	if now.Perm != "private" {
		t.Fatalf("新缓存对象 Perm 应为 private, got %q", now.Perm)
	}
	if now.Name != "Bob" {
		t.Fatalf("新缓存对象 Name 应为 Bob, got %q", now.Name)
	}
}

// 真实场景联动：controller/teacher.go 在 UpdatePerm 后会调 IsTeacher / GetTeacherName。
// 保证这两个读路径能看到更新后的数据。
func TestUpdateTeacherPerm_ReadPathReturnsFreshPerm(t *testing.T) {
	resetCache()
	teacherCache.Store(9, &entity.Teachers{Id: 9, Name: "Carol", Perm: "public"})

	UpdateTeacherPerm(9, "Carol", "private")

	// IsTeacher / GetTeacherName 在命中缓存时不会做 DB 调用，直接走缓存分支，可以放心调。
	perm, err := IsTeacher(nil, 9)
	if err != nil {
		t.Fatalf("IsTeacher err: %v", err)
	}
	if perm != "private" {
		t.Fatalf("IsTeacher 应看到 private, got %q", perm)
	}
	name, err := GetTeacherName(nil, 9)
	if err != nil {
		t.Fatalf("GetTeacherName err: %v", err)
	}
	if name != "Carol" {
		t.Fatalf("GetTeacherName 应看到 Carol, got %q", name)
	}
}
