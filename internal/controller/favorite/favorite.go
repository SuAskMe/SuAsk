package favorite

import (
	"context"
	v1 "suask/api/favorite/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
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
	if answersOutput != nil {
		for k, v := range answersOutput.AvatarsMap {
			dstId := QuestionList[idMap[k]].DstUserID
			if dstId != 0 {
				QuestionList[idMap[k]].AnswerAvatars = []string{consts.DefaultAvatarURL}
				continue
			}
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
			if len(URLs) == 0 {
				URLs = urls.URL
			} else {
				URLs = append(urls.URL, URLs[0])
			}
			QuestionList[idMap[k]].AnswerAvatars = URLs
		}
	}
	// 返回结果
	res = &v1.GetFavoritePageRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
	}
	return res, nil
}

func (c *cFavorite) Get(ctx context.Context, req *v1.GetFavoritePageReq) (res *v1.GetFavoritePageRes, err error) {
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

// func (c *cFavorite) GetKeyWords(ctx context.Context, req *v1.GetFavoriteSearchKeywordsReq) (res *v1.GetFavoriteSearchKeywordsRes, err error) {
// 	input := model.GetFavoriteKeywordsInput{}
// 	err = gconv.Scan(req, &input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	out, err := service.Favorite().GetKeyWord(ctx, &input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = gconv.Scan(out, &res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }

// func (c *cFavorite) GetByKeyWord(ctx context.Context, req *v1.GetFavoritePageByKeywordReq) (res *v1.GetFavoritePageByKeywordRes, err error) {
// 	data, err := GetFavoriteImpl(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = gconv.Scan(data, &res)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, err
// }

func (c *cFavorite) Favorite(ctx context.Context, req *v1.FavoriteReq) (res *v1.FavoriteRes, err error) {
	output, err := service.QuestionUtil().Favorite(ctx, &model.FavoriteInput{QuestionID: req.QuestionID})
	if err != nil {
		return nil, err
	}
	res = &v1.FavoriteRes{
		IsFavorite: output.IsFavorite,
	}
	return
}
