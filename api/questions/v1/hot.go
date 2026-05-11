package v1

import (
	"suask/internal/model"

	"github.com/gogf/gf/v2/frame/g"
)

// 热点问题列表（按浏览量排序）+ 全局搜索

type GetHotQuestionsReq struct {
	g.Meta    `path:"/questions/hot" method:"GET" tags:"Question" summary:"获取热点问题列表"`
	Page      int    `json:"page" v:"required|min:1"`
	TimeRange string `json:"time_range" v:"required|in:week,month,all" dc:"时间范围: week/month/all"`
	Keyword   string `json:"keyword" dc:"搜索关键词（按标题模糊匹配）"`
}

type GetHotQuestionsRes struct {
	QuestionList []model.TeacherQuestion `json:"question_list"`
	RemainPage   int                     `json:"remain_page"`
}
