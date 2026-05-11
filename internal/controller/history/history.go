package history

import (
	"context"
	v1 "suask/api/history/v1"
	"suask/internal/consts"
	qutil "suask/internal/logic/questions_util"
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
	// 获取图片
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	// 获取回答者头像
	answersOutput, err := service.QuestionUtil().GetAnswers(ctx, &model.GetAnswersInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	avatarsMap := map[int][]int{}
	if answersOutput != nil {
		avatarsMap = answersOutput.AvatarsMap
	}
	// 一次性批量查所有 file_id -> URL (#2 优化)
	allFileIDs := qutil.CollectFileIDs(imagesOutput.ImageMap, avatarsMap)
	urlMap, err := qutil.BatchGetFileURLs(ctx, allFileIDs)
	if err != nil {
		return nil, err
	}
	// 按原始顺序分发图片 URL
	imageURLs := qutil.ResolveImageURLs(urlMap, imagesOutput.ImageMap)
	for qid, urls := range imageURLs {
		QuestionList[idMap[qid]].ImageURLs = urls
	}
	// 按原始顺序分发头像 URL
	avatarURLs := qutil.ResolveAvatarURLs(urlMap, avatarsMap)
	for qid, urls := range avatarURLs {
		if QuestionList[idMap[qid]].DstUserID != 0 {
			QuestionList[idMap[qid]].AnswerAvatars = []string{consts.DefaultAvatarURL}
		} else {
			QuestionList[idMap[qid]].AnswerAvatars = urls
		}
	}
	// 返回结果
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
