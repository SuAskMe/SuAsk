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

func GetHistoryImpl(ctx context.Context, req interface{}) (res interface{}, err error) {
	baseInput := model.GetHistoryBaseInput{}
	err = gconv.Scan(req, &baseInput)
	if err != nil {
		return nil, err
	}
	baseOutput, err := service.History().GetBase(ctx, &baseInput)
	if err != nil {
		return nil, err
	}
	QuestionList := baseOutput.Questions
	idMap := baseOutput.IdMap
	dstUserIDMap := make(map[int]int, len(QuestionList))
	for _, question := range QuestionList {
		dstUserIDMap[question.ID] = question.DstUserID
	}
	assetsOutput, err := service.QuestionUtil().GetQuestionListAssets(ctx, &model.GetQuestionListAssetsInput{
		QuestionIDs:  baseOutput.QuestionIDs,
		DstUserIDMap: dstUserIDMap,
	})
	if err != nil {
		return
	}
	for questionId, urls := range assetsOutput.ImageURLMap {
		QuestionList[idMap[questionId]].ImageURLs = urls
	}
	for questionId, urls := range assetsOutput.AnswerAvatarMap {
		QuestionList[idMap[questionId]].AnswerAvatars = urls
	}
	res = &v1.GetHistoryPageRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
	}
	return res, nil
}

func (c *cHistory) Get(ctx context.Context, req *v1.GetHistoryPageReq) (res *v1.GetHistoryPageRes, err error) {
	data, err := GetHistoryImpl(ctx, req)
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(data, &res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *cHistory) GetKeyWords(ctx context.Context, req *v1.GetHistorySearchKeywordsReq) (res *v1.GetHistorySearchKeywordsRes, err error) {
	input := model.GetHistoryKeywordsInput{}
	err = gconv.Scan(req, &input)
	if err != nil {
		return nil, err
	}
	out, err := service.History().GetKeyWord(ctx, &input)
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(out, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *cHistory) GetByKeyWord(ctx context.Context, req *v1.GetHistoryPageByKeywordReq) (res *v1.GetHistoryPageByKeywordRes, err error) {
	data, err := GetHistoryImpl(ctx, req)
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(data, &res)
	if err != nil {
		return nil, err
	}
	return res, err
}
