package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cPublicQuestions struct{}

var PublicQuestions = cPublicQuestions{}

func (cPublicQuestions) Get(ctx context.Context, req *v1.GetPublicQuestionsReq) (res *v1.GetPublicQuestionsRes, err error) {
	input := model.GetPublicQuestionInput{}
	gconv.Scan(req, &input)
	ouput, err := service.PublicQuestion().Get(ctx, &input)
	res = &v1.GetPublicQuestionsRes{
		QuestionList: ouput.Questions,
		RemainPage:   ouput.RemainPage,
	}
	return
}
