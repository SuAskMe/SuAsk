package questions

import (
	"context"
	"strconv"
	v1 "suask/api/questions/v1"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cPublicQuestions struct{}

var PublicQuestions = cPublicQuestions{}

func (cPublicQuestions) Add(ctx context.Context, req *v1.AddQuestionReq) (res *v1.AddQuestionRes, err error) {
	questionInput := model.AddQuestionInput{}
	err = gconv.Struct(req, &questionInput)
	if err != nil {
		return nil, err
	}
	questionId, err := service.PublicQuestion().AddQuestion(ctx, &questionInput)
	if err != nil {
		return nil, err
	}
	if req.Files != nil {
		fileList := model.FileListAddInput{
			FileList: req.Files,
		}
		fileIdList, err := service.File().UploadFileList(ctx, fileList)
		if err != nil {
			return nil, err
		}
		attachment := model.AddAttachmentInput{
			QuestionId: questionId.ID,
			Type:       consts.QuestionFileType,
			FileId:     fileIdList.IdList,
		}
		_, err = service.Attachment().AddAttachments(ctx, attachment)
		if err != nil {
			return nil, err
		}
	}
	res = &v1.AddQuestionRes{Id: questionId.ID}
	return res, nil
}

func GetQuestionImpl(ctx context.Context, req interface{}) (res interface{}, err error) {
	baseInput := model.GetBaseInput{}
	gconv.Scan(req, &baseInput)
	baseOutput, err := service.PublicQuestion().GetBase(ctx, &baseInput)
	if err != nil {
		return
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
		urls, err_ := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
		if err_ != nil {
			return nil, err_
		}
		QuestionList[idMap[k]].AnswerNum = answersOutput.CountMap[k]
		QuestionList[idMap[k]].AnswerAvatars = urls.URL
	}
	// 返回结果
	res = &v1.GetPageRes{
		QuestionList: QuestionList,
		RemainPage:   baseOutput.RemainPage,
	}
	return
}

func (cPublicQuestions) Get(ctx context.Context, req *v1.GetPageReq) (res *v1.GetPageRes, err error) {
	res_, err := GetQuestionImpl(ctx, req)
	if err != nil {
		return
	}
	gconv.Scan(res_, &res)
	return
}

func (cPublicQuestions) GetKeywords(ctx context.Context, req *v1.GetSearchKeywordsReq) (res *v1.GetSearchKeywordsRes, err error) {
	input := model.GetKeywordsInput{}
	gconv.Scan(req, &input)
	ouput, err := service.PublicQuestion().GetKeyword(ctx, &input)
	gconv.Scan(ouput, &res)
	return
}

func (cPublicQuestions) GetByKeyword(ctx context.Context, req *v1.GetPageByKeywordReq) (res *v1.GetPageByKeywordRes, err error) {
	res_, err := GetQuestionImpl(ctx, req)
	if err != nil {
		return
	}
	gconv.Scan(res_, &res)
	return
}

func (cPublicQuestions) Favorite(ctx context.Context, req *v1.FavoriteReq) (res *v1.FavoriteRes, err error) {
	input := model.FavoriteInput{}
	gconv.Scan(req, &input)
	output, err := service.QuestionUtil().Favorite(ctx, &input)
	gconv.Scan(output, &res)
	return
}
