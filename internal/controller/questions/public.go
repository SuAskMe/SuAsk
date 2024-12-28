package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cPublicQuestions struct{}

var PublicQuestions = cPublicQuestions{}

func (cPublicQuestions) Add(ctx context.Context, req *v1.AddQuestionReq) (res *v1.AddQuestionRes, err error) {
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	//UserId := 2
	//if UserId == consts.DefaultUserId {
	//	return nil, fmt.Errorf("user not login")
	//}
	questionInput := model.AddQuestionInput{}
	err = gconv.Struct(req, &questionInput)
	questionInput.SrcUserID = UserId
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

	// 添加通知
	if req.DstUserId != 0 {
		_, err := service.Notification().Add(ctx, model.AddNotificationInput{
			UserId:     req.DstUserId,
			QuestionId: questionId.ID,
			Type:       consts.NewQuestion,
		})
		if err != nil {
			return nil, err
		}
	}
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

//func (cPublicQuestions) Favorite(ctx context.Context, req *v1.FavoriteReq) (res *v1.FavoriteRes, err error) {
//	input := model.FavoriteInput{}
//	gconv.Scan(req, &input)
//	output, err := service.QuestionUtil().Favorite(ctx, &input)
//	gconv.Scan(output, &res)
//	return
//}
