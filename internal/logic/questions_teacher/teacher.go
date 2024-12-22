package questions_teacher

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
)

type sTeacherQuestion struct{}

func sortByType(md **gdb.Model, sortType int) error {
	switch sortType {
	case consts.SortByTimeDsc:
		*md = (*md).Order("created_at DESC")
	case consts.SortByTimeAsc:
		*md = (*md).Order("created_at ASC")
	case consts.SortByViewsDsc:
		*md = (*md).Order("views DESC")
	case consts.SortByViewsAsc:
		*md = (*md).Order("views ASC")
	default:
		return fmt.Errorf("invalid sort type: %d", sortType)
	}
	return nil
}

// TruncateString 截断字符串：中文字符截断到 150 个字符，英文字符截断到 450 个字符
func TruncateString(s string) string {
	runes := []rune(s)
	length := 0
	for i, r := range runes {
		if r <= 0x7F {
			length++
		} else {
			length += 3
		}
		if length > 500 {
			return string(runes[:i]) + "..."
		}
	}
	return s
}

func (sTeacherQuestion) GetBase(ctx context.Context, input *model.GetBaseOfTeacherInput) (*model.GetBaseOfTeacherOutput, error) {
	md := dao.Questions.Ctx(ctx).Where("dst_user_id = ?", input.TeacherID)
	if input.Keyword != "" {
		md = md.Where("title LIKE?", "%"+input.Keyword+"%")
	}
	md = md.Page(input.Page, consts.NumOfQuestionsPerPage)
	err := sortByType(&md, input.SortType)
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
	// 获取问题ID列表 （官方的静态联表也是这么做的）
	qIDs := make([]int, len(q))
	for i, pq := range q {
		qIDs[i] = pq.Id
	}
	var fav []*custom.MyFavorites
	UserId := 1
	// UserId := gconv.Int(ctx.Value(consts.CtxId))
	md = dao.Favorites.Ctx(ctx).Where("question_id IN (?) AND user_id = ?", qIDs, UserId)
	err = md.Scan(&fav) // 再查favorites
	if err != nil {
		return nil, err
	}

	pqs := make([]model.TeacherQuestion, len(q)) // 用于存放最终结果
	idMap := make(map[int]int)                   // 用于快速查找问题ID对应的索引
	for i, pq := range q {
		idMap[pq.Id] = i
		pqs[i] = model.TeacherQuestion{
			ID:        pq.Id,
			Title:     pq.Title,
			Content:   TruncateString(pq.Contents),
			CreatedAt: pq.CreatedAt.TimestampMilli(),
			Views:     pq.Views,
		}
	}
	for _, f := range fav { // 填充IsFavorited字段
		pqs[idMap[f.QuestionId]].IsFavorited = true
	}

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
	md := dao.Questions.Ctx(ctx).Where("dst_user_id = ?", input.TeacherID)
	// fmt.Println(input.Keyword)
	err := sortByType(&md, input.SortType)
	if err != nil {
		return nil, err
	}
	words := make([]model.Keywords, 8)
	err = md.Where("title LIKE ?", "%"+input.Keyword+"%").Limit(8).Scan(&words)
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
