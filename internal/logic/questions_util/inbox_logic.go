package questions

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/middleware"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/model/do"
	"suask/internal/service"
	"suask/utility"

	"github.com/gogf/gf/v2/database/gdb"
)

// sTeacherQuestionSelf 实现老师收件箱的核心逻辑。
// 原来在 internal/logic/question_teacher_self/ 包里，
// 旧路由删除后挪到这里，由新的 /questions/inbox 接口复用。
type sTeacherQuestionSelf struct{}

const errInboxSearchFailed = "搜索失败，请稍后重试"

func (sTeacherQuestionSelf) GetQFMAll(ctx context.Context, input *model.GetQFMInput) (*model.GetQFMOutput, error) {
	relation := fmt.Sprintf("favorites.question_id = questions.id AND favorites.user_id = %d AND favorites.package = '%s'", input.TeacherId, consts.OnTop)
	page := input.Page
	if page < 1 {
		page = 1
	}
	limit := consts.MaxQuestionsPerPage
	offset := (page - 1) * limit

	if input.Tag == "pinned" {
		md := buildInboxQuestionQuery(ctx, input, relation).
			WhereNotNull("favorites.id")
		remain, err := md.Count()
		if err != nil {
			return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMAll: count pinned questions failed", "teacherId", input.TeacherId, "keyword", input.Keyword)
		}

		md = md.Fields("questions.*", "(favorites.id IS NOT NULL) AS is_pinned").
			Order("favorites.id DESC").
			Order("questions.created_at DESC").
			Limit(offset, limit)
		var q []*custom.Questions
		err = md.Scan(&q)
		if err != nil {
			return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMAll: query pinned questions failed", "teacherId", input.TeacherId, "page", page, "keyword", input.Keyword)
		}
		return buildQFMOutput(q, utility.CountRemainPage(remain, page)), nil
	}

	pinnedCountMd := buildInboxQuestionQuery(ctx, input, relation).
		WhereNotNull("favorites.id")
	pinnedCount, err := pinnedCountMd.Count()
	if err != nil {
		return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMAll: count pinned questions failed", "teacherId", input.TeacherId, "keyword", input.Keyword)
	}

	unpinnedCountMd := buildInboxQuestionQuery(ctx, input, relation).
		WhereNull("favorites.id")
	unpinnedCount, err := unpinnedCountMd.Count()
	if err != nil {
		return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMAll: count unpinned questions failed", "teacherId", input.TeacherId, "keyword", input.Keyword)
	}

	var q []*custom.Questions
	if offset < pinnedCount {
		pinnedLimit := limit
		if remainPinned := pinnedCount - offset; remainPinned < pinnedLimit {
			pinnedLimit = remainPinned
		}
		pinnedMd := buildInboxQuestionQuery(ctx, input, relation).
			WhereNotNull("favorites.id").
			Fields("questions.*", "(favorites.id IS NOT NULL) AS is_pinned").
			Order("favorites.id DESC").
			Order("questions.created_at DESC").
			Limit(offset, pinnedLimit)
		var pinnedQuestions []*custom.Questions
		err = pinnedMd.Scan(&pinnedQuestions)
		if err != nil {
			return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMAll: query pinned questions failed", "teacherId", input.TeacherId, "page", page, "keyword", input.Keyword)
		}
		q = append(q, pinnedQuestions...)
	}

	if len(q) < limit {
		unpinnedOffset := 0
		if offset > pinnedCount {
			unpinnedOffset = offset - pinnedCount
		}
		unpinnedLimit := limit - len(q)
		unpinnedMd := buildInboxQuestionQuery(ctx, input, relation).
			WhereNull("favorites.id").
			Fields("questions.*", "(favorites.id IS NOT NULL) AS is_pinned")
		err = utility.SortByType(&unpinnedMd, input.SortType)
		if err != nil {
			return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMAll: apply sort failed", "teacherId", input.TeacherId, "sortType", input.SortType)
		}
		unpinnedMd = unpinnedMd.Order("questions.id DESC").Limit(unpinnedOffset, unpinnedLimit)
		var unpinnedQuestions []*custom.Questions
		err = unpinnedMd.Scan(&unpinnedQuestions)
		if err != nil {
			return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMAll: query unpinned questions failed", "teacherId", input.TeacherId, "page", page, "keyword", input.Keyword)
		}
		q = append(q, unpinnedQuestions...)
	}

	remain := utility.CountRemainPage(pinnedCount+unpinnedCount, page)
	return buildQFMOutput(q, remain), nil
}

func buildInboxQuestionQuery(ctx context.Context, input *model.GetQFMInput, relation string) *gdb.Model {
	md := dao.Questions.Ctx(ctx)
	if input.Tag == "deleted" {
		md = md.Unscoped().WhereNotNull("questions.deleted_at")
	}
	md = md.LeftJoin("favorites", relation).
		Where(dao.Questions.Columns().DstUserId, input.TeacherId)

	switch input.Tag {
	case consts.Unanswered:
		md = md.Where("questions.reply_cnt", 0)
	case consts.Answered:
		md = md.WhereGT("questions.reply_cnt", 0)
	}

	if input.Keyword != "" {
		md = md.WhereLike("questions.title", "%"+input.Keyword+"%")
	}
	return md
}

func buildQFMOutput(q []*custom.Questions, remain int) *model.GetQFMOutput {
	qIDs := make([]int, len(q))
	pqs := make([]model.QFM, len(q))
	idMap := make(map[int]int)
	for i, pq := range q {
		idMap[pq.Id] = i
		qIDs[i] = pq.Id
		if pq.DeletedAt != nil {
			pqs[i].Tag = "已删除"
		} else if pq.ReplyCnt > 0 {
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

	return &model.GetQFMOutput{
		QuestionIDs: qIDs,
		IdMap:       idMap,
		Questions:   pqs,
		RemainPage:  remain,
	}
}

func (sTeacherQuestionSelf) GetQFMPinned(ctx context.Context, input *model.GetQFMInput) (*model.GetQFMOutput, error) {
	md := dao.Favorites.Ctx(ctx).Where(dao.Favorites.Columns().UserId, input.TeacherId)
	md = md.Where(dao.Favorites.Columns().Package, consts.OnTop)
	var fav []custom.MyFavorites
	err := md.Scan(&fav)
	if err != nil {
		return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMPinned: query pinned failed", "teacherId", input.TeacherId)
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
		return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.GetQFMPinned: query questions failed", "teacherId", input.TeacherId)
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
		return nil, middleware.SanitizeError(ctx, err, errInboxSearchFailed, "Inbox.GetKeyword: query keywords failed", "teacherId", input.TeacherId, "keyword", input.Keyword)
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
		return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.PinQFM: count pins failed", "teacherId", input.TeacherId, "questionId", input.QuestionId)
	}
	if cnt > 0 {
		_, err = md.Delete()
		if err != nil {
			return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.PinQFM: delete old pin failed", "teacherId", input.TeacherId, "questionId", input.QuestionId)
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
			return nil, middleware.SanitizeError(ctx, err, consts.ErrInternal, "Inbox.PinQFM: insert pin failed", "teacherId", input.TeacherId, "questionId", input.QuestionId)
		}
		return &model.PinQFMOutput{IsPinned: true}, nil
	}
}

func init() {
	service.RegisterTeacherQuestionSelf(&sTeacherQuestionSelf{})
}
