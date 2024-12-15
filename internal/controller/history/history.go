package history

import (
	"context"
	v1 "suask/api/history/v1"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cHistory struct{}

var History = cHistory{}

func (cHistory) Get(cxt context.Context, req *v1.LoadHistoryQuestionReq) (res *v1.LoadHistoryQuestionRes, err error) {
	// 声明与数据库交互的in结构体
	in := model.GetHistoryInput{}
	err = gconv.Scan(req, &in)
	if err != nil {
		return nil, err
	}
	// 使用接口中的注册接口方法调用login中的逻辑函数
	out, err := service.History().LoadHistoryInfo(cxt, &in)
	if err != nil {
		return nil, err
	}
	res = &v1.LoadHistoryQuestionRes{
		HistoryQuestionList: out.Question,
		RemainPage:          out.RemainPage,
	}
	return
}
