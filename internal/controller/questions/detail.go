package questions

import (
	"context"
	"fmt"
	"sync"

	v1 "suask/api/answer/v1"
	"suask/internal/consts"
	qutil "suask/internal/logic/questions_util"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/send_email"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type cQuestionDetail struct{}

var QuestionDetail = cQuestionDetail{}

func (cQuestionDetail) GetDetail(ctx context.Context, req *v1.GetDetailReq) (res *v1.GetDetailRes, err error) {
	qid := req.QuestionID
	userId := gconv.Int(ctx.Value(consts.CtxId))

	QBOutput, err := service.QuestionDetail().GetQuestionBase(ctx, &model.GetQuestionBaseInput{QuestionId: qid, UserId: userId})
	if err != nil {
		return nil, err
	}

	var (
		wg                sync.WaitGroup
		viewErr           error
		questionImageErr  error
		answerErr         error
		questionImageURLs []string
		ansOutput         *model.GetAnswerDetailOutput
	)
	wg.Add(3)
	go func() {
		defer wg.Done()
		_, viewErr = service.QuestionDetail().AddQuestionView(ctx, &model.AddViewInput{QuestionId: qid})
	}()
	go func() {
		defer wg.Done()
		urlMap, err := qutil.BatchGetFileURLs(ctx, QBOutput.ImageList)
		if err != nil {
			questionImageErr = err
			return
		}
		questionImageURLs = make([]string, 0, len(QBOutput.ImageList))
		for _, fid := range QBOutput.ImageList {
			if url, ok := urlMap[fid]; ok {
				questionImageURLs = append(questionImageURLs, url)
			}
		}
	}()
	go func() {
		defer wg.Done()
		ansOutput, answerErr = service.QuestionDetail().GetAnswers(ctx, &model.GetAnswerDetailInput{QuestionId: qid, DstUserId: QBOutput.DstUserId})
	}()
	wg.Wait()

	if viewErr != nil {
		return nil, viewErr
	}
	if questionImageErr != nil {
		return nil, questionImageErr
	}
	if answerErr != nil {
		return nil, answerErr
	}

	res = &v1.GetDetailRes{
		Question: model.QuestionBase{
			ID:         QBOutput.ID,
			Title:      QBOutput.Title,
			Content:    QBOutput.Content,
			Views:      QBOutput.Views + 1,
			CreatedAt:  QBOutput.CreatedAt,
			ImageURLs:  questionImageURLs,
			IsFavorite: QBOutput.IsFavorite,
		},
		CanReply: QBOutput.CanReply,
	}

	answerList := ansOutput.Answers
	IdMap := ansOutput.IdMap
	AvatarsMap := ansOutput.AvatarsMap
	ImageMap := ansOutput.ImageMap

	// 获取回答头像
	AvatarList := make([]int, 0, len(AvatarsMap))
	for k := range AvatarsMap {
		if k > 0 {
			AvatarList = append(AvatarList, k)
		}
	}

	// 收集所有需要查的 file_id：用户头像 + 回答图片，一次批查
	allFileIDs := make([]int, 0, len(AvatarList)+len(ImageMap)*4)
	allFileIDs = append(allFileIDs, AvatarList...)
	for _, fids := range ImageMap {
		allFileIDs = append(allFileIDs, fids...)
	}
	// 问题本身的图片已经在上面单独查过了（QBOutput.ImageList），这里只处理回答相关的
	urlMap, err := qutil.BatchGetFileURLs(ctx, allFileIDs)
	if err != nil {
		return nil, err
	}

	// 分发用户头像（file_id > 0 的）
	for fileId, ansIds := range AvatarsMap {
		if fileId > 0 {
			if url, ok := urlMap[fileId]; ok {
				for _, aid := range ansIds {
					answerList[IdMap[aid]].UserAvatar = url
				}
			}
		}
	}

	// 分发回答图片（按原始顺序）
	for answerId, fids := range ImageMap {
		urls := make([]string, 0, len(fids))
		for _, fid := range fids {
			if url, ok := urlMap[fid]; ok {
				urls = append(urls, url)
			}
		}
		answerList[IdMap[answerId]].ImageURLs = urls
	}
	res.Answers = answerList

	bgCtx := context.WithoutCancel(ctx)
	go func(userID, questionID int) {
		_, err := service.Notification().UpdateAoQ(bgCtx, model.UpdateAoQInput{UserID: userID, QuestionID: questionID})
		if err != nil {
			g.Log().Warningf(bgCtx, "更新问题通知已读状态失败: %v", err)
		}
	}(userId, qid)
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

	// 上传文件
	if req.Files != nil {
		fileList := model.FileListAddInput{
			UploaderId: UserId,
			FileList:   req.Files,
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

	// 添加通知（失败只记日志，不阻断回答已入库的事实）
	srcUserId, err := service.QuestionUtil().GetQuestionSrcUserId(ctx, req.QuestionId)
	if err != nil {
		g.Log().Warningf(ctx, "获取问题发起者失败: %v", err)
		return &v1.AddAnswerRes{Id: output.Id}, nil
	}
	// 给发帖的人通知有回答
	if srcUserId != consts.DefaultUserId && srcUserId != UserId {
		_, notifErr := service.Notification().Add(ctx, model.AddNotificationInput{
			UserId:     srcUserId,
			QuestionId: req.QuestionId,
			AnswerId:   output.Id,
			Type:       consts.NewAnswer,
		})
		if notifErr != nil {
			g.Log().Warningf(ctx, "添加回答通知失败: %v", notifErr)
		}
		emailErr := service.Notification().SendNoticeEmail(ctx, &model.SendNoticeEmailInput{
			To: srcUserId,
			Notice: &send_email.Notice{
				User:    "SuAsk用户",
				Type:    "新的回答",
				Content: input.Content,
				URL:     "https://suask.me/question-detail/" + gconv.String(req.QuestionId) + "#" + gconv.String(output.Id),
			},
		})
		if emailErr != nil {
			g.Log().Warningf(ctx, "发送回答邮件通知失败: %v", emailErr)
		}
	}

	// 如果是回复别人的回答
	if req.InReplyTo != nil {
		answer, ansErr := service.Answer().GetAnswerIDs(ctx, gconv.Int(req.InReplyTo))
		if ansErr != nil {
			g.Log().Warningf(ctx, "获取被回复的回答失败: %v", ansErr)
			return &v1.AddAnswerRes{Id: output.Id}, nil
		}
		// 回复不是默认用户或自己发的
		if answer.UserId != consts.DefaultUserId && answer.UserId != UserId {
			_, notifErr := service.Notification().Add(ctx, model.AddNotificationInput{
				UserId:     answer.UserId,
				AnswerId:   answer.Id,
				ReplyToId:  output.Id,
				QuestionId: answer.QuestionId,
				Type:       consts.NewReply,
			})
			if notifErr != nil {
				g.Log().Warningf(ctx, "添加回复通知失败: %v", notifErr)
			}
			emailErr := service.Notification().SendNoticeEmail(ctx, &model.SendNoticeEmailInput{
				To: answer.UserId,
				Notice: &send_email.Notice{
					User:    "SuAsk用户",
					Type:    "新的回复",
					Content: input.Content,
					URL:     "https://suask.me/question-detail/" + gconv.String(req.QuestionId) + "#" + gconv.String(output.Id),
				},
			})
			if emailErr != nil {
				g.Log().Warningf(ctx, "发送回复邮件通知失败: %v", emailErr)
			}
		}
	}

	res = &v1.AddAnswerRes{
		Id: output.Id,
	}
	return
}

func (cQuestionDetail) DeleteAnswer(ctx context.Context, req *v1.DeleteAnswerReq) (res *v1.DeleteAnswerRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录后操作")
	}
	if err := service.QuestionDetail().DeleteAnswer(ctx, req.ID, userId); err != nil {
		return nil, err
	}
	return &v1.DeleteAnswerRes{}, nil
}
