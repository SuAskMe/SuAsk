package questions_detail

import "suask/internal/model/custom"

// extractFileIDs 取出 attachments 查询结果里的 file_id 列表。
//
// 历史坑：之前这段是
//     ids := make([]int, len(imgList))          // len = N 的零值切片
//     for _, img := range imgList {
//         ids = append(ids, img.FileID)         // 再 append N 个 → 2N 长，前 N 个是 0
//     }
// 下游的 file.GetList(WhereIn ids) 会把 id=0 过滤掉，导致返回数量短于调用方预期，
// 进而让按 index 取 URL 的代码错位或越界。
//
// 独立出来便于单测，避免以后又被写回去。
func extractFileIDs(imgList []custom.Image) []int {
	out := make([]int, 0, len(imgList))
	for _, img := range imgList {
		out = append(out, img.FileID)
	}
	return out
}
