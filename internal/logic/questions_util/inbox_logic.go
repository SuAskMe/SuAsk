package questions

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

// sTeacherQuestionSelf 实现老师收件箱的核心逻辑。
// 原来在 internal/logic/question_teacher_self/ 包里，
// 旧路由删除后挪到这里，由新的 /questions/inbox 接口复用。
type sTeacherQuestionSelf struct{}

func (sTeacherQuestionSelf) GetQFMAll(ctx context.Context, input *model.GetQFMInput) (*model.GetQFMOutput, error) {
	relation := fmt.Sprintf("favorites.question_id = questions.id AND favorites.user_id = %d AND favorites.package = '%s'", input.TeacherId, consts.OnTop)
	md := dao.Questions.Ctx(ctx).
		LeftJoin("favorites", relation).
		Where(dao.Questions.Columns().DstUserId, input.TeacherId).
		WhereNull(dao.Questions.Columns().DeletedAt)

	switch input.Tag {
	case consts.Unanswered:
		md = md.Where("questions.reply_cnt", 0)
	case consts.Answered:
		md = md.WhereGT("questions.reply_cnt", 0)
	}

	if input.Keyword != "" {
		md = md.WhereLike("questions.title", "%"+input.Keyword+"%")
	}

	// 1. 先统计总数 (此时没有 Fields 和 Order，可生成正确的 COUNT(1) 语句)
	remain, err := md.Count()
	if err != nil {
		return nil, err
	}

	// 2. 注入 Fields、置顶以及选择的排序逻辑
	md = md.Fields("questions.*", "(favorites.id IS NOT NULL) AS is_pinned")
	md = md.Order("favorites.id DESC")
	err = utility.SortByType(&md, input.SortType)
	if err != nil {
		return nil, err
	}

	// 3. 应用分页进行列表查询
	md = md.Page(input.Page, consts.MaxQuestionsPerPage)
	var q []*custom.Questions
	err = md.Scan(&q)
	if err != nil {
		return nil, err
	}
	remain = utility.CountRemainPage(remain, input.Page)

	qIDs := make([]int, len(q))
	pqs := make([]model.QFM, len(q))
	idMap := make(map[int]int)
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
		pqs[i].IsPinned = pq.IsPinned
	}

	output := model.GetQFMOutput{
		QuestionIDs: qIDs,
		IdMap:       idMap,
		Questions:   pqs,
		RemainPage:  remain,
	}
	return &output, nil
}

func (sTeacherQuestionSelf) GetQFMPinned(ctx context.Context, input *model.GetQFMInput) (*model.GetQFMOutput, error) {
	md := dao.Favorites.Ctx(ctx).Where(dao.Favorites.Columns().UserId, input.TeacherId)
	md = md.Where(dao.Favorites.Columns().Package, consts.OnTop)
	var fav []custom.MyFavorites
	err := md.Scan(&fav)
	if err != nil {
		return nil, err
	}
	qIDs := make([]int, len(fav))
	for i, f := range fav {
		qIDs[i] = f.QuestionId
	}
	if len(qIDs) == 0 {
		return &model.GetQFMOutput{Questions: []model.QFM{}}, nil
	}

	md = dao.Questions.Ctx(ctx).WhereIn(dao.Questions.Columns().Id, qIDs).Where("deleted_at IS NULL")
	var q []*custom.Questions
	err = md.Scan(&q)
	if err != nil {
		return nil, err
	}
	pqs := make([]model.QFM, len(q))
	idMap := make(map[int]int)
	for i, pq := range q {
		idMap[pq.Id] = i
		pqs[i].ID = pq.Id
		pqs[i].Title = pq.Title
		pqs[i].Content = utility.TruncateString(pq.Contents)
		pqs[i].CreatedAt = pq.CreatedAt.TimestampMilli()
		pqs[i].Views = pq.Views
		pqs[i].IsPinned = true
		if pq.ReplyCnt > 0 {
			pqs[i].Tag = consts.Answered
		} else {
			pqs[i].Tag = consts.Unanswered
		}
	}
	newQIDs := make([]int, len(pqs))
	for i, pq := range pqs {
		newQIDs[i] = pq.ID
	}
	return &model.GetQFMOutput{
		QuestionIDs: newQIDs,
		IdMap:       idMap,
		Questions:   pqs,
	}, nil
}

func (sTeacherQuestionSelf) GetKeyword(ctx context.Context, input *model.GetQFMKeywordsInput) (*model.GetKeywordsOutput, error) {
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().DstUserId, input.TeacherId).Where("deleted_at IS NULL")
	md = md.WhereLike(dao.Questions.Columns().Title, "%"+input.Keyword+"%").Limit(8)
	words := make([]model.Keyword, consts.MaxKeywordsPerReq)
	err := md.Scan(&words)
	if err != nil {
		return nil, nil
	}
	output := &model.GetKeywordsOutput{}
	output.Words = words
	return output, nil
}

func (sTeacherQuestionSelf) PinQFM(ctx context.Context, input *model.PinQFMInput) (*model.PinQFMOutput, error) {
	md := dao.Favorites.Ctx(ctx).Where(dao.Favorites.Columns().UserId, input.TeacherId)
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
		return &model.PinQFMOutput{IsPinned: false}, nil
	} else {
		md = dao.Favorites.Ctx(ctx)
		_, err = md.Insert(do.Favorites{
			UserId:     input.TeacherId,
			QuestionId: input.QuestionId,
			Package:    consts.OnTop,
		})
		if err != nil {
			return nil, err
		}
		return &model.PinQFMOutput{IsPinned: true}, nil
	}
}

func init() {
	service.RegisterTeacherQuestionSelf(&sTeacherQuestionSelf{})
}
