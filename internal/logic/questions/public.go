package public

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

type sPublicQuestion struct{}

// var keywordCacheMode = gdb.CacheOption{
// 	Duration: time.Minute * 5,
// 	Name:     "public_keywords_cache",
// 	Force:    false,
// }

// var keywordClearCache = gdb.CacheOption{
// 	Duration: -1,
// 	Name:     "public_keywords_cache",
// 	Force:    false,
// }

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

func (sPublicQuestion) Get(ctx context.Context, input *model.GetInput) (*model.GetOutput, error) {
	md := dao.Questions.Ctx(ctx).WhereNull("dst_user_id")
	if input.Keyword != "" {
		md = md.Where("title LIKE?", "%"+input.Keyword+"%")
	}
	md = md.Page(input.Page, consts.NumOfQuestionsPerPage)
	md = md.WithAll()
	err := sortByType(&md, input.SortType)
	if err != nil {
		return nil, err
	}
	var q []*custom.PublicQuestions
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
	md = dao.Favorites.Ctx(ctx).Where("question_id IN (?) AND user_id = ?", qIDs, input.UserID)
	err = md.Scan(&fav) // 再查favorites
	if err != nil {
		return nil, err
	}

	pqs := make([]model.PublicQuestion, len(q)) // 用于存放最终结果
	idMap := make(map[int]int)                  // 用于快速查找问题ID对应的索引
	for i, pq := range q {
		idMap[pq.Id] = i
		pqs[i] = model.PublicQuestion{
			ID:            pq.Id,
			Title:         pq.Title,
			Content:       TruncateString(pq.Contents),
			CreatedAt:     pq.CreatedAt.TimestampMilli(),
			Views:         pq.Views,
			ImageURLs:     nil,
			IsFavorited:   false,
			AnswerNum:     len(pq.Answers),
			AnswerAvatars: nil,
		}
	}
	for _, f := range fav { // 填充IsFavorited字段
		pqs[idMap[f.QuestionId]].IsFavorited = true
	}

	output := model.GetOutput{
		Questions:  pqs,
		RemainPage: remain,
	}
	return &output, nil
}

func (sPublicQuestion) GetKeyword(ctx context.Context, input *model.GetKeywordsInput) (*model.GetKeywordsOutput, error) {
	// md := dao.Questions.Ctx(ctx).Cache(keywordCacheMode).WhereNull("dst_user_id")
	md := dao.Questions.Ctx(ctx).WhereNull("dst_user_id")
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

func (sPublicQuestion) Favorite(ctx context.Context, input *model.FavoriteInput) (*model.FavoriteOutput, error) {
	md := dao.Favorites.Ctx(ctx)
	cnt, err := md.Where("user_id = ? AND question_id = ?", input.UserID, input.QuestionID).Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		_, err = md.Delete("user_id = ? AND question_id = ?", input.UserID, input.QuestionID)
		if err != nil {
			return nil, err
		}
		return &model.FavoriteOutput{
			IsFavorited: false,
		}, nil
	} else {
		_, err = md.Insert(g.Map{
			"user_id":     input.UserID,
			"question_id": input.QuestionID,
		})
		if err != nil {
			return nil, err
		}
		return &model.FavoriteOutput{
			IsFavorited: true,
		}, nil
	}
}

func init() {
	service.RegisterPublicQuestion(New())
}

func New() *sPublicQuestion {
	return &sPublicQuestion{}
}
