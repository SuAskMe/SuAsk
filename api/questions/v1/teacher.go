package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

type GetPageOfTeacherReq struct {
	g.Meta `path:"/questions/teacher" method:"GET" tags:"Question" summary:"获取老师问题列表"`
	GetPageBase
	TeacherID int `v:"required|min:1" json:"teacher_id"`
}

type GetPageOfTeacherRes struct {
	QuestionList []model.TeacherQuestion `json:"question_list"`
	RemainPage   int                     `json:"remain_page"`
}

type GetSearchKeywordsOfTeacherReq struct {
	g.Meta    `path:"/questions/teacher/keywords" method:"GET" tags:"Question" summary:"搜索老师问题"`
	Keyword   string `v:"required|length:2,100" json:"keyword"`
	SortType  int    `v:"required|min:0|max:3" json:"sort_type"`
	TeacherID int    `v:"required|min:1" json:"teacher_id"`
}

type GetSearchKeywordsOfTeacherRes struct {
	Words []struct {
		Value string `json:"value"`
	} `json:"words"`
}

type GetPageByKeywordOfTeacherReq struct {
	g.Meta    `path:"/questions/teacher/search" method:"GET" tags:"Question" summary:"根据关键字获取老师问题列表"`
	Keyword   string `v:"length:2,100" json:"keyword"`
	TeacherID int    `v:"required|min:1" json:"teacher_id"`
	GetPageBase
}

type GetPageByKeywordOfTeacherRes struct {
	QuestionList []model.TeacherQuestion `json:"question_list"`
	RemainPage   int                     `json:"remain_page"`
}

//type FavoriteOfTeacherReq struct {
//	g.Meta     `path:"/questions/teacher/favorite" method:"post" tags:"Teacher" summary:"收藏老师问题"`
//	QuestionID int `v:"required|min:1" json:"question_id"`
//}
//
//type FavoriteOfTeacherRes struct {
//	IsFavorited bool `json:"is_favorited"`
//}
