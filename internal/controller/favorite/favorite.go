package favorite

import (
	"context"
	v1 "suask/api/favorite/v1"
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
	for questionId, users := range assetsOutput.AnswerUserMap {
		QuestionList[idMap[questionId]].AnswerUsers = users
	}
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
