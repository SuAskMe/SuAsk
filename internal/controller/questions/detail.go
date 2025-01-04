package questions

import (
	"context"
	"fmt"
	v1 "suask/api/answer/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/util/gconv"
)

type cQuestionDetail struct{}

var QuestionDetail = cQuestionDetail{}

func (cQuestionDetail) GetDetail(ctx context.Context, req *v1.GetDetailReq) (res *v1.GetDetailRes, err error) {
	qid := req.QuestionID
	userId := gconv.Int(ctx.Value(consts.CtxId))

	// 获取问题
	QBOutput, err := service.QuestionDetail().GetQuestionBase(ctx, &model.GetQuestionBaseInput{QuestionId: qid, UserId: userId})
	if err != nil {
		return nil, err
	}
	// 增加浏览量
	_, err = service.QuestionDetail().AddQuestionView(ctx, &model.AddViewInput{QuestionId: qid})
	if err != nil {
		return nil, err
	}
	// 获取问题图片列表
	fileList, err := service.File().GetList(ctx, model.FileListGetInput{IdList: QBOutput.ImageList})
	if err != nil {
		return nil, err
	}
	res = &v1.GetDetailRes{ // 填充返回结果
		Question: model.QuestionBase{
			ID:         QBOutput.ID,
			Title:      QBOutput.Title,
			Content:    QBOutput.Content,
			Views:      QBOutput.Views,
			CreatedAt:  QBOutput.CreatedAt,
			ImageURLs:  fileList.URL,
			IsFavorite: QBOutput.IsFavorite,
		},
		CanReply: QBOutput.CanReply,
	}
	// 获取回答列表
	ansOutput, err := service.QuestionDetail().GetAnswers(ctx, &model.GetAnswerDetailInput{QuestionId: qid, DstUserId: QBOutput.DstUserId})
	if err != nil {
		return nil, err
	}

	answerList := ansOutput.Answers
	IdMap := ansOutput.IdMap
	AvatarsMap := ansOutput.AvatarsMap
	ImageMap := ansOutput.ImageMap

	// 获取回答头像
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
	// 获取回答图片
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

func (cQuestionDetail) Upvote(ctx context.Context, req *v1.UpvoteReq) (res *v1.UpvoteRes, err error) {
	input := model.UpvoteInput{}
	gconv.Scan(req, &input)
	output, err := service.QuestionDetail().AddAnswerUpvote(ctx, &input)
	if err != nil {
		return
	}
	gconv.Scan(output, &res)
	return
}

func (cQuestionDetail) AddAnswer(ctx context.Context, req *v1.AddAnswerReq) (res *v1.AddAnswerRes, err error) {
	if req.Content == "" {
		return nil, fmt.Errorf("content is empty")
	}
	input := model.AddAnswerInput{}
	err = gconv.Scan(req, &input)
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	input.UserId = UserId
	if err != nil {
		return nil, err
	}
	output, err := service.QuestionDetail().ReplyQuestion(ctx, &input)
	if err != nil {
		return
	}

	// 更新回答数
	replyCntOut, err := service.QuestionDetail().AddReplyCnt(ctx, &model.AddReplyCntInput{QuestionId: input.QuestionId})
	if err != nil {
		return
	}

	// 记录头像
	if replyCntOut.ReplyCnt <= consts.MaxAvatarsPerQuestion {
		_, err = service.QuestionDetail().BuildRelation(ctx, &model.BuildRelationInput{QuestionId: input.QuestionId})
		if err != nil {
			return nil, err
		}
	}

	// 上传文件
	if req.Files != nil {
		fileList := model.FileListAddInput{
			FileList: req.Files,
		}
		fileIdList, err := service.File().UploadFileList(ctx, fileList)
		if err != nil {
			return nil, err
		}
		attachment := model.AddAttachmentInput{
			AnswerId: output.Id,
			Type:     consts.QuestionFileType,
			FileId:   fileIdList.IdList,
		}
		_, err = service.Attachment().AddAttachments(ctx, attachment)
		if err != nil {
			return nil, err
		}
	}

	// 添加通知
	srcUserId, err := service.QuestionUtil().GetQuestionSrcUserId(ctx, req.QuestionId)
	if err != nil {
		return nil, err
	}
	// 给发帖的人通知有回答
	if srcUserId != consts.DefaultUserId && srcUserId != UserId {
		_, err = service.Notification().Add(ctx, model.AddNotificationInput{
			UserId:     srcUserId,
			QuestionId: req.QuestionId,
			AnswerId:   output.Id,
			Type:       consts.NewAnswer,
		})
	}

	// 如果是回复别人的回答
	if req.InReplyTo != nil {
		answer, err := service.Answer().GetAnswerIDs(ctx, gconv.Int(req.InReplyTo))
		if err != nil {
			return nil, err
		}
		// 回复不是默认用户或自己发的
		if answer.UserId != consts.DefaultUserId && answer.UserId != UserId {
			_, err := service.Notification().Add(ctx, model.AddNotificationInput{
				UserId:     answer.UserId,
				AnswerId:   answer.Id,
				ReplyToId:  output.Id,
				QuestionId: answer.QuestionId,
				Type:       consts.NewReply,
			})
			if err != nil {
				return nil, err
			}
		}
	}

	if err != nil {
		return nil, err
	}
	res = &v1.AddAnswerRes{
		Id: output.Id,
	}
	return
}
