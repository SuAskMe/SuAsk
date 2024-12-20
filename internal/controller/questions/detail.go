package questions

import (
	"context"
	v1 "suask/api/questions/v1"
	"suask/internal/model"
	"suask/internal/service"
)

type cQuestionDetail struct{}

var QuestionDetail = cQuestionDetail{}

func (cQuestionDetail) GetDetail(ctx context.Context, req *v1.GetDetailReq) (res *v1.GetDetailRes, err error) {
	qid := req.QuestionID
	questionBaseOutput, err := service.QuestionDetail().GetQuestionBase(ctx, &model.GetQuestionBaseInput{QuestionId: qid})
	if err != nil {
		return nil, err
	}
	res = &v1.GetDetailRes{
		Question: *questionBaseOutput.Question,
		CanReply: questionBaseOutput.CanReply,
	}
	answersOutput, err := service.QuestionDetail().GetAnswers(ctx, &model.GetAnswerDetailInput{QuestionId: qid})
	if err != nil {
		return nil, err
	}

	answerList := answersOutput.Answers
	IdMap := answersOutput.IdMap
	AvatarsMap := answersOutput.AvatarsMap // 将用户头像id映射到answer id列表
	ImageMap := answersOutput.ImageMap

	AvatarList := make([]int, 0, len(AvatarsMap))
	for k := range AvatarsMap {
		if k != 0 {
			AvatarList = append(AvatarList, k)
		}
	}
	avatarUrls, err := service.File().GetList(ctx, model.FileListGetInput{IdList: AvatarList})
	if err != nil {
		return nil, err
	}
	for i, url := range avatarUrls.URL {
		IdList := AvatarsMap[avatarUrls.FileId[i]]
		for _, id := range IdList {
			answerList[IdMap[id]].UserAvatar = url
		}
	}

	for k, v := range ImageMap {
		url, err := service.File().GetList(ctx, model.FileListGetInput{IdList: v})
		if err != nil {
			return nil, err
		}
		answerList[IdMap[k]].ImageURLs = url.URL
	}
	res.Answers = answerList
	return
}
