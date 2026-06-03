package questions_teacher

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"
	"suask/utility"
)

type sTeacherQuestion struct{}

func (sTeacherQuestion) GetBase(ctx context.Context, input *model.GetBaseOfTeacherInput) (*model.GetBaseOfTeacherOutput, error) {
	relation := fmt.Sprintf("favorites.question_id = questions.id AND favorites.user_id = %d AND favorites.package = '%s'", input.TeacherID, consts.OnTop)
	md := dao.Questions.Ctx(ctx).
		LeftJoin("favorites", relation).
		Where(dao.Questions.Columns().DstUserId, input.TeacherID).
		WhereNull(dao.Questions.Columns().DeletedAt).
		WhereGT(dao.Questions.Columns().ReplyCnt, 0)

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

	// 3. 再应用分页进行列表查询
	md = md.Page(input.Page, consts.MaxQuestionsPerPage)
	var q []*custom.Questions
	err = md.Scan(&q)
	if err != nil {
		return nil, err
	}
	// 计算剩余页数
	remain = utility.CountRemainPage(remain, input.Page)
	// 获取问题ID列表
	qIDs := make([]int, len(q))
	for i, pq := range q {
		qIDs[i] = pq.Id
	}

	pqs := make([]model.TeacherQuestion, len(q)) // 用于存放最终结果
	idMap := make(map[int]int)                   // 用于快速查找问题ID对应的索引
	for i, pq := range q {
		idMap[pq.Id] = i
		pqs[i] = model.TeacherQuestion{
			ID:        pq.Id,
			Title:     pq.Title,
			Content:   utility.TruncateString(pq.Contents),
			CreatedAt: pq.CreatedAt.TimestampMilli(),
			Views:     pq.Views,
			IsPinned:  pq.IsPinned,
		}
	}
	//for _, f := range fav { // 填充IsFavorited字段
	//	pqs[idMap[f.QuestionId]].IsFavorited = true
	//}

	output := model.GetBaseOfTeacherOutput{
		QuestionIDs: qIDs,
		IdMap:       idMap,
		Questions:   pqs,
		RemainPage:  remain,
	}
	return &output, nil
}

func (sTeacherQuestion) GetKeyword(ctx context.Context, input *model.GetKeywordsOfTeacherInput) (*model.GetKeywordsOutput, error) {
	// md := dao.Questions.Ctx(ctx).Cache(keywordCacheMode).WhereNull("dst_user_id")
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().DstUserId, input.TeacherID)
	md = md.WhereGT(dao.Questions.Columns().ReplyCnt, 0)
	// fmt.Println(input.Keyword)
	words := make([]model.Keyword, consts.MaxKeywordsPerReq)
	err := md.WhereLike(dao.Questions.Columns().Title, "%"+input.Keyword+"%").Limit(8).Scan(&words)
	if err != nil {
		return nil, err
	}
	output := &model.GetKeywordsOutput{}
	output.Words = words
	return output, nil
}

func init() {
	service.RegisterTeacherQuestion(New())
}

func New() *sTeacherQuestion {
	return &sTeacherQuestion{}
}
