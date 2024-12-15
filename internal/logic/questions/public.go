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

func (sPublicQuestion) Get(ctx context.Context, input *model.GetInput) (*model.GetOutput, error) {
	var q []*custom.PublicQuestions
	md := dao.Questions.Ctx(ctx).WhereNull("dst_user_id")
	if input.Keyword != "" {
		md = md.Where("title LIKE?", "%"+input.Keyword+"%")
	}
	qList := md.Page(input.Page, consts.NumOfQuestionsPerPage)
	qList = qList.WithAll()
	qList = qList.Where(custom.UserUpvotes{UserID: input.UserID}).Where(custom.UserFavorites{UserID: input.UserID})
	err := sortByType(&qList, input.SortType)
	if err != nil {
		return nil, err
	}
	err = qList.Scan(&q)
	if err != nil {
		return nil, err
	}
	pqs := make([]model.PublicQuestion, len(q))
	for i, pq := range q {
		pqs[i] = model.PublicQuestion{
			ID:            pq.Id,
			Title:         pq.Title,
			Content:       pq.Contents,
			CreatedAt:     pq.CreatedAt,
			Views:         pq.Views,
			ImageURLs:     nil,
			IsFavorited:   len(pq.IsFavorited) == 1,
			IsUpvoted:     len(pq.IsUpvoted) == 1,
			AnswerNum:     len(pq.Answers),
			AnswerAvatars: nil,
		}
	}
	remain, err := md.Count()
	if err != nil {
		return nil, err
	}
	remainNum := remain - consts.NumOfQuestionsPerPage*input.Page
	remain = remainNum / consts.NumOfQuestionsPerPage
	if remainNum%consts.NumOfQuestionsPerPage > 0 {
		remain += 1
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
	words := make([]model.Keywords, 10)
	err = md.Where("title LIKE ?", "%"+input.Keyword+"%").Limit(10).Scan(&words)
	if err != nil {
		return nil, err
	}
	output := &model.GetKeywordsOutput{}
	output.Words = words
	return output, nil
}

func (sPublicQuestion) Favorite(ctx context.Context, input *model.FavoriteInput) error {
	md := dao.Favorites.Ctx(ctx)
	cnt, err := md.Where("user_id = ? AND question_id = ?", input.UserID, input.QuestionID).Count()
	if err != nil {
		return err
	}
	if cnt > 0 {
		return fmt.Errorf("already favorited")
	}
	_, err = md.Insert(g.Map{
		"user_id":     input.UserID,
		"question_id": input.QuestionID,
	})
	return err
}

func init() {
	service.RegisterPublicQuestion(New())
}

func New() *sPublicQuestion {
	return &sPublicQuestion{}
}
