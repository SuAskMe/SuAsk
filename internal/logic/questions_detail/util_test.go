package questions_detail

import (
	"suask/internal/model/custom"
	"testing"
)

func TestExtractFileIDs_NoLeadingZeros(t *testing.T) {
	imgs := []custom.Image{
		{QuestionId: 1, FileID: 10},
		{QuestionId: 1, FileID: 20},
		{QuestionId: 1, FileID: 30},
	}
	got := extractFileIDs(imgs)

	if len(got) != len(imgs) {
		t.Fatalf("长度应为 %d，实际 %d: %v", len(imgs), len(got), got)
	}
	for i, id := range got {
		if id == 0 {
			t.Fatalf("第 %d 个元素不应是 0（调用方会用作 WhereIn 的 file_id），got=%v", i, got)
		}
	}
	want := []int{10, 20, 30}
	for i := range got {
		if got[i] != want[i] {
			t.Fatalf("顺序错乱：got=%v, want=%v", got, want)
		}
	}
}

func TestExtractFileIDs_EmptyInput(t *testing.T) {
	if got := extractFileIDs(nil); len(got) != 0 {
		t.Fatalf("空输入应返回空切片，got=%v", got)
	}
	if got := extractFileIDs([]custom.Image{}); len(got) != 0 {
		t.Fatalf("空切片输入应返回空切片，got=%v", got)
	}
}
