package history

import (
	"context"
	v1 "suask/api/history/v1"
	"suask/internal/model"
	"suask/internal/service"
)

type cHistory struct{}

var History = cHistory{}

func (cHistory) Get(cxt context.Context, req *v1.LoadHistoryQuestionReq) (res *v1.LoadHistoryQuestionRes, err error) {
	//userId = gconv.Int(ctx.Value(consts.CtxId))
	userId := 1
	in := model.GetHistoryInput{
		UserId: userId,
		Page:   req.Page,
	}

	// 使用接口中的注册接口方法调用login中的逻辑函数
	out, err := service.HistoryOperation().LoadHistoryInfo(cxt, &in)
	if err != nil {
		return nil, err
	}
	res = &v1.LoadHistoryQuestionRes{
		HistoryQuestionList: out.Question,
		Total:               out.Total,
		Size:                out.Size,
		PageNum:             out.PageNum,
		RemainPage:          out.RemainPage,
	}
	return
}
