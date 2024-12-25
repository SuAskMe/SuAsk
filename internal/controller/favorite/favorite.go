package favorite

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	v2 "suask/api/favorite/v2"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"
)

type cFavorite struct{}

var Favorite cFavorite

func GetFavoriteImpl(ctx context.Context, req interface{}) (res interface{}, err error) {
	baseInput := model.GetFavoriteBaseInput{}
	err = gconv.Scan(req, &baseInput)
	if err != nil {
		return nil, err
	}
	baseOutput, err := service.Favorite().GetBase(ctx, &baseInput)
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
		QuestionList[idMap[k]].AnswerNum = answersOutput.CountMap[k]
		QuestionList[idMap[k]].AnswerAvatars = URLs
	}
	// 返回结果
	res = &v2.GetPageRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
	}
	return res, nil
}

func (c *cFavorite) Get(ctx context.Context, req *v2.GetPageReq) (res *v2.GetPageRes, err error) {
	data, err := GetFavoriteImpl(ctx, req)
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(data, &res)
	if err != nil {
		return nil, err
	}
	return res, err
}

func (c *cFavorite) GetKeyWords(ctx context.Context, req *v2.GetSearchKeywordsReq) (res *v2.GetSearchKeywordsRes, err error) {
	input := model.GetFavoriteKeywordsInput{}
	err = gconv.Scan(req, &input)
	if err != nil {
		return nil, err
	}
	out, err := service.Favorite().GetKeyWord(ctx, &input)
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(out, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (c *cFavorite) GetByKeyWord(ctx context.Context, req *v2.GetPageByKeywordReq) (res *v2.GetPageByKeywordRes, err error) {
	data, err := GetFavoriteImpl(ctx, req)
	if err != nil {
		return nil, err
	}
	err = gconv.Scan(data, &res)
	if err != nil {
		return nil, err
	}
	return res, err
}
