package favorite

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

type sFavorite struct{}

func (s *sFavorite) GetBase(ctx context.Context, in *model.GetFavoriteBaseInput) (out *model.GetFavoriteBaseOutput, err error) {
	md := dao.Favorites.Ctx(ctx)
	//userId := 1
	userId := gconv.Int(ctx.Value(consts.CtxId))
	md = md.Where(dao.Favorites.Columns().UserId, userId)
	md = md.Page(in.Page, consts.NumOfQuestionsPerPage)
	err = utility.SortByType(&md, in.SortType)
	if err != nil {
		return nil, err
	}
	var remain int
	var f []*model.Favorite
	err = md.ScanAndCount(&f, &remain, true)
	if err != nil {
		return nil, err
	}
	remainNum := remain - consts.NumOfQuestionsPerPage*in.Page
	remain = remainNum / consts.NumOfQuestionsPerPage
	if remainNum%consts.NumOfQuestionsPerPage > 0 {
		remain += 1
	}
	qIDs := make([]int, len(f))
	for i, favorite := range f {
		qIDs[i] = favorite.QuestionId
	}
	var q []*custom.Questions
	md = dao.Questions.Ctx(ctx).WhereIn(dao.Questions.Columns().Id, qIDs)
	md = md.WhereLike(dao.Questions.Columns().Title, "%"+in.Keyword+"%")
	err = md.Scan(&q)
	if err != nil {
		return nil, err
	}
	pqs := make([]model.PublicQuestion, len(q))
	qMap := make(map[int]*custom.Questions, len(q))
	for _, question := range q {
		qMap[question.Id] = question
	}
	idMap := make(map[int]int)
	for i, id := range qIDs {
		if q, ok := qMap[id]; ok {
			idMap[q.Id] = i
			pqs[i] = model.PublicQuestion{
				ID:         q.Id,
				Title:      q.Title,
				Content:    utility.TruncateString(q.Contents),
				CreatedAt:  f[i].CreatedAt.TimestampMilli(),
				Views:      q.Views,
				AnswerNum:  q.ReplyCnt,
				IsFavorite: true,
			}
		}
	}

	output := &model.GetFavoriteBaseOutput{
		QuestionIDs: qIDs,
		Questions:   pqs,
		IdMap:       idMap,
		RemainPage:  remain,
	}
	return output, err
}

func (s *sFavorite) GetKeyWord(ctx context.Context, in *model.GetFavoriteKeywordsInput) (out *model.GetFavoriteKeywordsOutput, err error) {
	md := dao.Favorites.Ctx(ctx)
	//userId := 1
	userId := gconv.Int(ctx.Value(consts.CtxId))
	md = md.Where(dao.Favorites.Columns().UserId, userId)
	err = utility.SortByType(&md, in.SortType)
	if err != nil {
		return nil, err
	}
	var f []*model.Favorite
	var count int
	err = md.ScanAndCount(&f, &count, true)
	if err != nil {
		return nil, err
	}
	qIDs := make([]int, count)
	for i, favorite := range f {
		qIDs[i] = favorite.QuestionId
	}
	words := make([]model.Keyword, consts.NumOfKeywordsPerReq)
	md = dao.Questions.Ctx(ctx)
	err = md.WhereIn(dao.Questions.Columns().Id, qIDs).WhereLike(dao.Questions.Columns().Title, "%"+in.Keyword+"%").Limit(8).Scan(&words)
	if err != nil {
		return nil, err
	}
	out = &model.GetFavoriteKeywordsOutput{
		Words: words,
	}
	return out, nil
}

func init() {
	service.RegisterFavorite(New())
}

func New() *sFavorite {
	return &sFavorite{}
}
