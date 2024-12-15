package star

import (
	"context"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
)

type sStar struct{}

func (s *sStar) Star(ctx context.Context, in model.StarQuestionInPut) (out model.StarQuestionOutPut, err error) {
	id := 1 // 先写死了
	star := model.StarQuestionOutPut{}
	err = dao.Favorites.Ctx(ctx).With(model.QuestionInfo{}).Where(do.Favorites{UserId: id}).Scan(&star)

	return
}
