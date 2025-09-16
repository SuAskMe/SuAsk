package answer

import (
	"context"
	"suask/internal/dao"
	"suask/internal/model/entity"
	"suask/internal/service"
)

type sAnswer struct{}

func (s *sAnswer) GetAnswerIDs(ctx context.Context, answerId int) (out entity.Answers, err error) {
	out = entity.Answers{}
	err = dao.Answers.Ctx(ctx).Where(dao.Answers.Columns().Id, answerId).
		Fields("id, user_id, question_id").Scan(&out)
	if err != nil {
		return entity.Answers{}, err
	}
	return out, nil
}

func init() {
	service.RegisterAnswer(New())
}

func New() *sAnswer {
	return &sAnswer{}
}
