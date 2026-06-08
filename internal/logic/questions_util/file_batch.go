package questions

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model/entity"
	files "suask/utility/files"

	"github.com/gogf/gf/v2/frame/g"
)

// BatchGetFileURLs 一次性批量查询 file_id → URL 的映射。
// 调用方拿到 map 后按自己持有的 id 列表顺序取值，保证图片顺序不乱。
//
// 这是 #2 性能优化的核心：把原来 for 循环里逐个 File.GetList 的 N 次 SQL
// 合并成 1 次 WhereIn。
func BatchGetFileURLs(ctx context.Context, fileIDs []int) (map[int]string, error) {
	fileIDs = CollectUniqueFileIDs(fileIDs)
	if len(fileIDs) == 0 {
		return map[int]string{}, nil
	}
	var fileList []entity.Files
	err := dao.Files.Ctx(ctx).WhereIn(dao.Files.Columns().Id, fileIDs).Scan(&fileList)
	if err != nil {
		return nil, err
	}
	urlMap := make(map[int]string, len(fileList))
	for _, f := range fileList {
		url, err := files.GetURL(f.Hash, f.Name)
		if err != nil {
			g.Log().Warningf(ctx, "BatchGetFileURLs: file %d GetURL err: %v", f.Id, err)
			continue
		}
		urlMap[f.Id] = url
	}
	return urlMap, nil
}

// ResolveImageURLs 把 imageMap[questionId][]fileId 转成 imageURLs[questionId][]string，
// 保持每个问题内的图片顺序。
func ResolveImageURLs(urlMap map[int]string, imageMap map[int][]int) map[int][]string {
	result := make(map[int][]string, len(imageMap))
	for qid, fids := range imageMap {
		urls := make([]string, 0, len(fids))
		for _, fid := range fids {
			if url, ok := urlMap[fid]; ok {
				urls = append(urls, url)
			}
		}
		result[qid] = urls
	}
	return result
}

// ResolveAvatarURLs 把 avatarsMap[questionId][]avatarFileId 转成 URL 列表，
// avatarFileId=0 的用默认头像。保持顺序。
func ResolveAvatarURLs(urlMap map[int]string, avatarsMap map[int][]int) map[int][]string {
	result := make(map[int][]string, len(avatarsMap))
	for qid, fids := range avatarsMap {
		urls := make([]string, 0, len(fids))
		for _, fid := range fids {
			if fid == 0 {
				urls = append(urls, consts.DefaultAvatarURL)
			} else if url, ok := urlMap[fid]; ok {
				urls = append(urls, url)
			}
		}
		result[qid] = urls
	}
	return result
}

func CollectUniqueFileIDs(fileIDs []int) []int {
	if len(fileIDs) == 0 {
		return []int{}
	}
	seen := make(map[int]struct{}, len(fileIDs))
	ids := make([]int, 0, len(fileIDs))
	for _, fid := range fileIDs {
		if fid == 0 {
			continue
		}
		if _, ok := seen[fid]; ok {
			continue
		}
		seen[fid] = struct{}{}
		ids = append(ids, fid)
	}
	return ids
}

// CollectFileIDs 从 imageMap + avatarsMap 收集所有非零 file_id 并去重。
func CollectFileIDs(imageMap map[int][]int, avatarsMap map[int][]int) []int {
	ids := make([]int, 0)
	for _, fids := range imageMap {
		ids = append(ids, fids...)
	}
	for _, fids := range avatarsMap {
		ids = append(ids, fids...)
	}
	return CollectUniqueFileIDs(ids)
}
