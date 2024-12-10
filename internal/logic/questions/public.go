package public

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"
)

type sPublicQuestion struct{}

func (p *sPublicQuestion) Get(ctx context.Context, input *model.GetPublicQuestionInput) (*model.GetPublicQuestionOutput, error) {
	var q []*custom.PublicQuestions
	md := dao.Questions.Ctx(ctx).WhereNull("dst_user_id")
	qList := md.Page(input.Page, consts.NumOfQuestionsPerPage)
	qList = qList.WithAll()
	qList = qList.Where(custom.UserUpvotes{UserID: input.UserID}).Where(custom.UserFavorites{UserID: input.UserID})
	switch input.SortType {
	case consts.SortByTimeDsc:
		qList = qList.Order("created_at DESC")
	case consts.SortByTimeAsc:
		qList = qList.Order("created_at ASC")
	case consts.SortByViewsDsc:
		qList = qList.Order("views DESC")
	case consts.SortByViewsAsc:
		qList = qList.Order("views ASC")
	default:
		return nil, fmt.Errorf("invalid sort type: %d", input.SortType)
	}
	err := qList.Scan(&q)
	if err != nil {
		return nil, err
	}
	pqs := make([]model.PublicQuestion, len(q))
	for i, pq := range q {
		pqs[i] = model.PublicQuestion{
			ID:            pq.Id,
			Content:       &pq.Contents,
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
	output := model.GetPublicQuestionOutput{
		Questions:  pqs,
		RemainPage: remain,
	}
	return &output, nil
}

func init() {
	service.RegisterPublicQuestion(New())
}

func New() *sPublicQuestion {
	return &sPublicQuestion{}
}
