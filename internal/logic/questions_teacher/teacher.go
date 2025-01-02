package questions_teacher

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"
	"suask/utility"
)

type sTeacherQuestion struct{}

func (sTeacherQuestion) GetBase(ctx context.Context, input *model.GetBaseOfTeacherInput) (*model.GetBaseOfTeacherOutput, error) {
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().DstUserId, input.TeacherID)
	if input.Keyword != "" {
		md = md.Where("match(title) against (? in boolean mode)", input.Keyword)
	}
	md = md.Page(input.Page, consts.MaxQuestionsPerPage)
	err := utility.SortByType(&md, input.SortType)
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
	remain = utility.CountRemainPage(remain, input.Page)
	// 获取问题ID列表 （官方的静态联表也是这么做的）
	qIDs := make([]int, len(q))
	for i, pq := range q {
		qIDs[i] = pq.Id
	}
	//var fav []*custom.MyFavorites
	//UserId := 1
	//// UserId := gconv.Int(ctx.Value(consts.CtxId))
	//md = dao.Favorites.Ctx(ctx).WhereIn(dao.Favorites.Columns().QuestionId, qIDs).WhereIn(dao.Favorites.Columns().UserId, UserId)
	//err = md.Scan(&fav) // 再查favorites
	//if err != nil {
	//	return nil, err
	//}

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
	// fmt.Println(input.Keyword)
	err := utility.SortByType(&md, input.SortType)
	if err != nil {
		return nil, err
	}
	words := make([]model.Keyword, consts.MaxKeywordsPerReq)
	err = md.Where("match(title) against (? in boolean mode)", input.Keyword).Limit(8).Scan(&words)
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
