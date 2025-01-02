package utility

import (
	"fmt"
	"suask/internal/consts"

	"github.com/gogf/gf/v2/database/gdb"
)

func SortByType(md **gdb.Model, sortType int) error {
	switch sortType {
	case consts.SortByTimeDsc:
		*md = (*md).Order("created_at DESC")
	case consts.SortByTimeAsc:
		*md = (*md).Order("created_at ASC")
	case consts.SortByViewsDsc:
		*md = (*md).Order("views DESC")
	case consts.SortByViewsAsc:
		*md = (*md).Order("views ASC")
	default:
		return fmt.Errorf("invalid sort type: %d", sortType)
	}
	return nil
}

// TruncateString 截断字符串：中文字符截断到 150 个字符，英文字符截断到 450 个字符
func TruncateString(s string) string {
	runes := []rune(s)
	length := 0
	for i, r := range runes {
		if r <= 0x7F {
			length++
		} else {
			length += 3
		}
		if length > 500 {
			return string(runes[:i]) + "..."
		}
	}
	return s
}

func CountRemainPage(remain, page int) int {
	remainNum := remain - consts.MaxQuestionsPerPage*page
	remain = remainNum / consts.MaxQuestionsPerPage
	if remainNum%consts.MaxQuestionsPerPage > 0 {
		remain += 1
	}
	return remain
}

func AddUnique[T comparable](slice []T, value T) []T {
	uniqueMap := make(map[T]bool)
	for _, v := range slice {
		uniqueMap[v] = true
	}
	if !uniqueMap[value] {
		slice = append(slice, value)
	}

	return slice
}
