package history

import (
	"context"
	v1 "suask/api/history/v1"
	"suask/internal/consts"
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
	for k, v := range imagesOutput.ImageMap {
		urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
		if err_ != nil {
			return nil, err_
		}
		QuestionList[idMap[k]].ImageURLs = urls.URL
	}
	// 获取回答数
	answersOutput, err := service.PublicQuestion().GetAnswers(ctx, &model.GetAnswersInput{QuestionIDs: baseOutput.QuestionIDs})
	if err != nil {
		return
	}
	for k, v := range answersOutput.AvatarsMap {
		idList := make([]int, 0, len(v))
		URLs := make([]string, 0, len(v))
		for _, u := range v {
			if u != 0 {
				idList = append(idList, u)
			} else {
				URLs = append(URLs, consts.DefaultAvatarURL)
			}
		}
		urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: idList})
		if err_ != nil {
			return nil, err_
		}
		URLs = append(URLs, urls.URL...)
		QuestionList[idMap[k]].AnswerAvatars = URLs
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
