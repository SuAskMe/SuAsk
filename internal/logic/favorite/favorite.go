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
	userId := gconv.Int(ctx.Value(consts.CtxId))

	// 关联 questions 表进行过滤与排序，确保能按 questions.created_at / questions.views 正常运行
	md := dao.Favorites.Ctx(ctx).
		LeftJoin("questions", "questions.id = favorites.question_id").
		Where(dao.Favorites.Columns().UserId, userId).
		Where(dao.Favorites.Columns().Package, "default").
		WhereNull("questions.deleted_at")

	// 1. 先在干净的 Model 上统计总数 (无 Fields, Order, Page 干扰)
	remain, err := md.Count()
	if err != nil {
		return nil, err
	}

	// 2. 注入特定的字段 (只读取 favorites.* 以便 Scan 映射)
	md = md.Fields("favorites.*")

	// 3. 注入排序逻辑
	err = utility.SortByType(&md, in.SortType)
	if err != nil {
		return nil, err
	}

	// 4. 应用分页进行列表查询
	md = md.Page(in.Page, consts.MaxQuestionsPerPage)
	var f []*model.Favorite
	err = md.Scan(&f)
	if err != nil {
		return nil, err
	}

	remain = utility.CountRemainPage(remain, in.Page)

	qIDs := make([]int, len(f))
	for i, favorite := range f {
		qIDs[i] = favorite.QuestionId
	}

	if len(qIDs) == 0 {
		output := &model.GetFavoriteBaseOutput{
			QuestionIDs: qIDs,
			Questions:   []*model.FavoriteQuestion{},
			IdMap:       map[int]int{},
			RemainPage:  remain,
		}
		return output, nil
	}

	var q []custom.Questions
	md = dao.Questions.Ctx(ctx).WhereIn(dao.Questions.Columns().Id, qIDs)
	err = md.Scan(&q)
	// fmt.Println("q:", q)
	if err != nil {
		return nil, err
	}
	pqs := make([]*model.FavoriteQuestion, 0, len(q))
	qMap := make(map[int]*custom.Questions, len(q))
	for _, question := range q {
		if question.DstUserId != 0 && question.ReplyCnt == 0 { // 严防通过favorite获取到未回复的对老师的提问
			continue
		}
		qMap[question.Id] = &question
	}
	idMap := make(map[int]int)
	for i, id := range qIDs {
		if q, ok := qMap[id]; ok {
			fq := &model.FavoriteQuestion{
				ID:         q.Id,
				Title:      q.Title,
				Content:    utility.TruncateString(q.Contents),
				CreatedAt:  f[i].CreatedAt.TimestampMilli(),
				Views:      q.Views,
				AnswerNum:  q.ReplyCnt,
				DstUserID:  q.DstUserId,
				IsFavorite: true,
			}
			pqs = append(pqs, fq)
			idMap[q.Id] = len(pqs) - 1
		}
	}
	qIDs = make([]int, len(pqs))
	for i, pq := range pqs {
		qIDs[i] = pq.ID
	}
	// fmt.Println("qIDs:", qIDs)
	output := &model.GetFavoriteBaseOutput{
		QuestionIDs: qIDs,
		Questions:   pqs,
		IdMap:       idMap,
		RemainPage:  remain,
	}
	return output, err
}

// func (s *sFavorite) GetKeyWord(ctx context.Context, in *model.GetFavoriteKeywordsInput) (out *model.GetFavoriteKeywordsOutput, err error) {
// 	md := dao.Favorites.Ctx(ctx)
// 	//userId := 1
// 	userId := gconv.Int(ctx.Value(consts.CtxId))
// 	md = md.Where(dao.Favorites.Columns().UserId, userId)
// 	err = utility.SortByType(&md, in.SortType)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var f []*model.Favorite
// 	var count int
// 	err = md.ScanAndCount(&f, &count, true)
// 	if err != nil {
// 		return nil, err
// 	}
// 	qIDs := make([]int, count)
// 	for i, favorite := range f {
// 		qIDs[i] = favorite.QuestionId
// 	}
// 	words := make([]model.Keyword, consts.NumOfKeywordsPerReq)
// 	md = dao.Questions.Ctx(ctx)
// 	err = md.WhereIn(dao.Questions.Columns().Id, qIDs).WhereLike(dao.Questions.Columns().Title, "%"+in.Keyword+"%").Limit(8).Scan(&words)
// 	if err != nil {
// 		return nil, err
// 	}
// 	out = &model.GetFavoriteKeywordsOutput{
// 		Words: words,
// 	}
// 	return out, nil
// }

func init() {
	service.RegisterFavorite(New())
}

func New() *sFavorite {
	return &sFavorite{}
}
