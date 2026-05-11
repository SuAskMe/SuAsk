package favorite

import (
	"context"
	v1 "suask/api/favorite/v1"
	"suask/internal/consts"
	qutil "suask/internal/logic/questions_util"
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
		dstId := QuestionList[idMap[qid]].DstUserID
		if dstId != 0 {
			QuestionList[idMap[qid]].AnswerAvatars = []string{consts.DefaultAvatarURL}
		} else {
			QuestionList[idMap[qid]].AnswerAvatars = urls
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
