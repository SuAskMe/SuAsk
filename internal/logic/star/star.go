package star

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

type MyStarQuestion struct {
	g.Meta     `orm:"table:favorites"`
	QuestionId int         `json:"question_id"`
	StarAt     *gtime.Time `json:"star_at"`

	Questions *QuestionInfo `orm:"with:id=question_id" json:"questions"`
}

type sStar struct{}

func (s *sStar) GetStar(ctx context.Context, id int) (out model.StarQuestionOutPut, err error) {
	var star []*MyStarQuestion
	err = dao.Favorites.Ctx(ctx).With(QuestionInfo{}).Where(do.Favorites{UserId: id}).Scan(&star)

	var list []model.StarQuestion
	for _, v := range star {
		var question model.StarQuestion
		question.ID = v.QuestionId
		question.Title = v.Questions.Title
		question.Contents = v.Questions.Contents
		question.Views = v.Questions.Views
		question.StarAt = v.StarAt
		list = append(list, question)
	}
	StarList := model.StarQuestionOutPut{StarQuestionList: list}
	return StarList, err
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

func init() {
	service.RegisterStar(New())
}

func New() *sStar {
	return &sStar{}
}
