package star

import (
	"context"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
)

type sStar struct{}

func (s *sStar) GetStar(ctx context.Context, id int) (out model.StarQuestionOutPut, err error) {
	star := model.StarQuestionOutPut{}
	err = dao.Favorites.Ctx(ctx).With(model.QuestionInfo{}).Where(do.Favorites{UserId: id}).Scan(&star)

	return
}
