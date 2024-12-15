package star

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
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

func (s *sStar) DeleteStar(ctx context.Context, in model.DeleteStarInput) (out model.DeleteStarOutput, err error) {
	_, err = g.Model("Favorites").Ctx(ctx).Where("id", in.Id).Delete()
	if err != nil {
		return model.DeleteStarOutput{}, err
	}
	res := model.DeleteStarOutput{
		String: "取消收藏成功",
	}
	return res, err
}
