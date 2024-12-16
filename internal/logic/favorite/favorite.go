package favorite

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
)

type QuestionInfo struct {
	g.Meta   `orm:"table:questions"`
	ID       int    `json:"id" dc:"问题ID"`
	Title    string `json:"title" dc:"标题"`
	Contents string `json:"contents" dc:"问题内容"`
	Views    int    `json:"views" dc:"浏览量"`
}

type MyFavoriteQuestion struct {
	g.Meta     `orm:"table:favorite"`
	QuestionId int         `json:"question_id"`
	FavoriteAt *gtime.Time `json:"favorite_at"`

	Questions *QuestionInfo `orm:"with:id=question_id" json:"questions"`
}

type sFavorite struct{}

func (s *sFavorite) GetFavorite(ctx context.Context, id int) (out model.FavoriteQuestionOutPut, err error) {
	var favorite []*MyFavoriteQuestion
	err = dao.Favorites.Ctx(ctx).With(QuestionInfo{}).Where(do.Favorites{UserId: id}).Scan(&favorite)

	var list []model.FavoriteQuestion
	for _, v := range favorite {
		var question model.FavoriteQuestion
		question.ID = v.QuestionId
		question.Title = v.Questions.Title
		question.Contents = v.Questions.Contents
		question.Views = v.Questions.Views
		question.FavoriteAt = v.FavoriteAt
		list = append(list, question)
	}
	FavoriteList := model.FavoriteQuestionOutPut{FavoriteQuestionList: list}
	return FavoriteList, err
}

func (s *sFavorite) DeleteFavorite(ctx context.Context, in model.DeleteFavoriteInput) (out model.DeleteFavoriteOutput, err error) {
	_, err = g.Model("favorite").Ctx(ctx).Where("id", in.Id).Delete()
	if err != nil {
		return model.DeleteFavoriteOutput{}, err
	}
	res := model.DeleteFavoriteOutput{
		String: "取消收藏成功",
	}
	return res, err
}

func init() {
	service.RegisterFavorite(New())
}

func New() *sFavorite {
	return &sFavorite{}
}
