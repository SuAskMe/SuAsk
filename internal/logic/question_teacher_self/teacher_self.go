package questions_teacher_self

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/model/do"
	"suask/internal/service"
	"suask/utility"
)

type sTeacherQuestionSelf struct{}

func Validate(ctx context.Context) error {
	// Tid := gconv.Int(ctx.Value(consts.CtxId))
	Tid := 2
	md := dao.Users.Ctx(ctx).Where(dao.Users.Columns().Id, Tid).Fields(dao.Users.Columns().Role)
	var role string
	err := md.Scan(&role)
	if err != nil {
		return err
	}
	if role != consts.TEACHER {
		return fmt.Errorf("user is not a teacher")
	}
	return nil
}

func (sTeacherQuestionSelf) GetQFMAll(ctx context.Context, input *model.GetQFMInput) (*model.GetQFMOutput, error) {
	err := Validate(ctx)
	if err != nil {
		return nil, err
	}
	Tid := 2
	// Tid := gconv.Int(ctx.Value(consts.CtxId))
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().DstUserId, Tid)
	switch input.Tag {
	case consts.Unanswered:
		md = md.Where(dao.Questions.Columns().ReplyCnt, 0)
	case consts.Answered:
		md = md.Where(dao.Questions.Columns().ReplyCnt, ">", 0)
	}
	if input.Keyword != "" {
		md = md.WhereLike(dao.Questions.Columns().Title, "%"+input.Keyword+"%")
	}
	md = md.Page(input.Page, consts.NumOfQuestionsPerPage)
	err = utility.SortByType(&md, input.SortType)
	if err != nil {
		return nil, err
	}

	var q []*custom.Questions
	var remain int
	err = md.ScanAndCount(&q, &remain, true) // 先查不包含favorites的结果
	if err != nil {
		return nil, err
	}
	// 计算剩余页数
	remainNum := remain - consts.NumOfQuestionsPerPage*input.Page
	remain = remainNum / consts.NumOfQuestionsPerPage
	if remainNum%consts.NumOfQuestionsPerPage > 0 {
		remain += 1
	}

	qIDs := make([]int, len(q))      // 获取问题ID列表 （官方的静态联表也是这么做的）
	pqs := make([]model.QFM, len(q)) // 用于存放最终结果
	idMap := make(map[int]int)       // 用于快速查找问题ID对应的索引
	for i, pq := range q {
		idMap[pq.Id] = i
		qIDs[i] = pq.Id
		if pq.ReplyCnt > 0 {
			pqs[i].Tag = consts.Answered
		} else {
			pqs[i].Tag = consts.Unanswered
		}
		pqs[i].ID = pq.Id
		pqs[i].Title = pq.Title
		pqs[i].Content = utility.TruncateString(pq.Contents)
		pqs[i].CreatedAt = pq.CreatedAt.TimestampMilli()
		pqs[i].Views = pq.Views
	}

	var fav []*custom.MyFavorites
	md = dao.Favorites.Ctx(ctx).WhereIn(dao.Favorites.Columns().QuestionId, qIDs)
	md = md.Where(dao.Favorites.Columns().UserId, Tid)
	md = md.Where(dao.Favorites.Columns().Package, consts.OnTop)
	err = md.Scan(&fav) // 再查favorites
	if err != nil {
		return nil, err
	}
	for _, f := range fav { // 填充IsFavorited字段
		// pqs[idMap[f.QuestionId]].Tag = consts.OnTop
		pqs[idMap[f.QuestionId]].IsPinned = true
	}
	output := model.GetQFMOutput{
		QuestionIDs: qIDs,
		IdMap:       idMap,
		Questions:   pqs,
		RemainPage:  remain,
	}
	return &output, nil
}

func (sTeacherQuestionSelf) GetQFMPinned(ctx context.Context, _ *model.GetQFMInput) (*model.GetQFMOutput, error) {
	// 置顶问题，不分页，直接返回全部
	err := Validate(ctx)
	if err != nil {
		return nil, err
	}
	Tid := 2
	// Tid := gconv.Int(ctx.Value(consts.CtxId))
	md := dao.Favorites.Ctx(ctx).Where(dao.Favorites.Columns().UserId, Tid)
	md = md.Where(dao.Favorites.Columns().Package, consts.OnTop)
	var fav []custom.MyFavorites
	err = md.Scan(&fav)
	if err != nil {
		return nil, err
	}
	qIDs := make([]int, len(fav))
	for i, f := range fav {
		qIDs[i] = f.QuestionId
	}

	md = dao.Questions.Ctx(ctx).WhereIn(dao.Questions.Columns().Id, qIDs)

	var q []*custom.Questions
	err = md.Scan(&q) // 先查不包含favorites的结果
	if err != nil {
		return nil, err
	}
	pqs := make([]model.QFM, len(q)) // 用于存放最终结果
	idMap := make(map[int]int)       // 用于快速查找问题ID对应的索引
	for i, pq := range q {
		idMap[pq.Id] = i
		qIDs[i] = pq.Id
		pqs[i].Tag = consts.OnTop
		pqs[i].ID = pq.Id
		pqs[i].Title = pq.Title
		pqs[i].Content = utility.TruncateString(pq.Contents)
		pqs[i].CreatedAt = pq.CreatedAt.TimestampMilli()
		pqs[i].Views = pq.Views
	}
	output := model.GetQFMOutput{
		QuestionIDs: qIDs,
		IdMap:       idMap,
		Questions:   pqs,
	}
	return &output, nil
}

func (sTeacherQuestionSelf) GetKeyword(ctx context.Context, input *model.GetQFMKeywordsInput) (*model.GetKeywordsOutput, error) {
	err := Validate(ctx)
	if err != nil {
		return nil, err
	}
	Tid := 2
	// Tid := gconv.Int(ctx.Value(consts.CtxId))
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().DstUserId, Tid)
	md = md.WhereLike(dao.Questions.Columns().Title, "%"+input.Keyword+"%").Limit(8)
	err = utility.SortByType(&md, input.SortType)
	if err != nil {
		return nil, err
	}
	words := make([]model.Keywords, 8)
	err = md.Scan(&words)
	if err != nil {
		return nil, err
	}
	output := &model.GetKeywordsOutput{}
	output.Words = words
	return output, nil
}

func (sTeacherQuestionSelf) PinQFM(ctx context.Context, input *model.PinQFMInput) (*model.PinQFMOutput, error) {
	err := Validate(ctx)
	if err != nil {
		return nil, err
	}
	Tid := 2
	// Tid := gconv.Int(ctx.Value(consts.CtxId))
	md := dao.Favorites.Ctx(ctx).Where(dao.Favorites.Columns().UserId, Tid)
	md = md.Where(dao.Favorites.Columns().QuestionId, input.QuestionId)
	md = md.Where(dao.Favorites.Columns().Package, consts.OnTop)
	cnt, err := md.Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		_, err = md.Delete()
		if err != nil {
			return nil, err
		}
		return &model.PinQFMOutput{
			IsPinned: false,
		}, nil
	} else {
		md = dao.Favorites.Ctx(ctx)
		_, err = md.Insert(do.Favorites{
			UserId:     Tid,
			QuestionId: input.QuestionId,
			Package:    consts.OnTop,
		})
		if err != nil {
			return nil, err
		}
		return &model.PinQFMOutput{
			IsPinned: true,
		}, nil
	}
}

func init() {
	service.RegisterTeacherQuestionSelf(New())
}

func New() *sTeacherQuestionSelf {
	return &sTeacherQuestionSelf{}
}
