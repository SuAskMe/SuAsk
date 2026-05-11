package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	"suask/internal/dao"
	qutil "suask/internal/logic/questions_util"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"
	"suask/utility"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

type cHotQuestion struct{}

var HotQuestion = cHotQuestion{}

func (c *cHotQuestion) Get(ctx context.Context, req *v1.GetHotQuestionsReq) (res *v1.GetHotQuestionsRes, err error) {
	md := dao.Questions.Ctx(ctx).Where("deleted_at IS NULL")
	// 只展示有回答的问题
	md = md.WhereGT(dao.Questions.Columns().ReplyCnt, 0)
	// 只展示非私密问题
	md = md.Where(dao.Questions.Columns().IsPrivate, 0)

	// 时间范围筛选
	switch req.TimeRange {
	case "week":
		md = md.WhereGTE(dao.Questions.Columns().CreatedAt, gtime.New(time.Now().AddDate(0, 0, -7)))
	case "month":
		md = md.WhereGTE(dao.Questions.Columns().CreatedAt, gtime.New(time.Now().AddDate(0, -1, 0)))
		// "all" 不加时间条件
	}

	// 关键词搜索
	if req.Keyword != "" {
		md = md.WhereLike(dao.Questions.Columns().Title, "%"+req.Keyword+"%")
	}

	// 按浏览量降序排序
	md = md.OrderDesc(dao.Questions.Columns().Views)
	md = md.Page(req.Page, consts.MaxQuestionsPerPage)

	var q []*custom.Questions
	var total int
	err = md.ScanAndCount(&q, &total, false)
	if err != nil {
		return nil, err
	}

	remain := utility.CountRemainPage(total, req.Page)

	qIDs := make([]int, len(q))
	pqs := make([]model.TeacherQuestion, len(q))
	idMap := make(map[int]int)
	for i, pq := range q {
		qIDs[i] = pq.Id
		idMap[pq.Id] = i
		pqs[i] = model.TeacherQuestion{
			ID:        pq.Id,
			Title:     pq.Title,
			Content:   utility.TruncateString(pq.Contents),
			CreatedAt: pq.CreatedAt.TimestampMilli(),
			Views:     pq.Views,
		}
	}

	// 获取图片
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: qIDs})
	if err != nil {
		return nil, err
	}
	if imagesOutput != nil && len(imagesOutput.ImageMap) > 0 {
		allFileIDs := qutil.CollectFileIDs(imagesOutput.ImageMap, nil)
		urlMap, err_ := qutil.BatchGetFileURLs(ctx, allFileIDs)
		if err_ != nil {
			return nil, err_
		}
		imageURLs := qutil.ResolveImageURLs(urlMap, imagesOutput.ImageMap)
		for qid, urls := range imageURLs {
			pqs[idMap[qid]].ImageURLs = urls
		}
	}

	res = &v1.GetHotQuestionsRes{
		QuestionList: pqs,
		RemainPage:   remain,
	}
	return
}
