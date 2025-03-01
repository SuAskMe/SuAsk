package history

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"
	"suask/utility"

	"github.com/gogf/gf/v2/util/gconv"
)

type sHistory struct{}

func (s *sHistory) GetBase(ctx context.Context, in *model.GetHistoryBaseInput) (out *model.GetHistoryBaseOutput, err error) {
	//userId := 1
	userId := gconv.Int(ctx.Value(consts.CtxId))

	md := dao.Questions.Ctx(ctx)
	md = md.Where(dao.Questions.Columns().SrcUserId, userId)
	if in.Keyword != "" {
		md = md.Where("match(title) against (? in boolean mode)", in.Keyword)
	} else {
		err = utility.SortByType(&md, in.SortType)
		if err != nil {
			return nil, err
		}
	}

	md = md.Page(in.Page, consts.MaxQuestionsPerPage)
	var remain int
	var q []*custom.Questions
	err = md.ScanAndCount(&q, &remain, false)
	if err != nil {
		return nil, err
	}

	remain = utility.CountRemainPage(remain, in.Page)

	qIDs := make([]int, len(q))
	for i, question := range q {
		qIDs[i] = question.Id
	}
	var fav []*custom.MyFavorites
	md = dao.Favorites.Ctx(ctx).WhereIn(dao.Favorites.Columns().QuestionId, qIDs).Where(dao.Favorites.Columns().UserId, userId)
	err = md.Scan(&fav)
	if err != nil {
		return nil, err
	}
	pqs := make([]model.HistoryQuestion, len(q))
	idMap := make(map[int]int)
	for i, pq := range q {
		idMap[pq.Id] = i
		pqs[i] = model.HistoryQuestion{
			ID:        pq.Id,
			Title:     pq.Title,
			Content:   utility.TruncateString(pq.Contents),
			CreatedAt: pq.CreatedAt.TimestampMilli(),
			Views:     pq.Views,
			AnswerNum: pq.ReplyCnt,
			DstUserID: pq.DstUserId,
		}
	}
	for _, f := range fav {
		pqs[idMap[f.QuestionId]].IsFavorite = true
	}
	output := model.GetHistoryBaseOutput{
		QuestionIDs: qIDs,
		IdMap:       idMap,
		Questions:   pqs,
		RemainPage:  remain,
	}
	return &output, err
}

func (s *sHistory) GetKeyWord(ctx context.Context, in *model.GetHistoryKeywordsInput) (out *model.GetHistoryKeywordsOutput, err error) {
	//userId := 1
	userId := gconv.Int(ctx.Value(consts.CtxId))

	md := dao.Questions.Ctx(ctx)
	md = md.Where(dao.Questions.Columns().SrcUserId, userId)
	words := make([]model.Keyword, consts.MaxKeywordsPerReq)
	err = md.Where("match(title) against (? in boolean mode)", in.Keyword).Limit(8).Scan(&words)
	if err != nil {
		return nil, err
	}
	output := &model.GetHistoryKeywordsOutput{}
	output.Words = words
	return output, nil
}

func init() {
	service.RegisterHistory(New())
}

func New() *sHistory {
	return &sHistory{}
}
