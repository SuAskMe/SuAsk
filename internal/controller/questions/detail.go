package questions

import (
	"context"
	"fmt"
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
			Views:      QBOutput.Views + 1,
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
	TeacherAvatarList := make([]int, 0)
	for k := range AvatarsMap {
		if k > 0 {
			AvatarList = append(AvatarList, k)
		} else if k < 0 {
			TeacherAvatarList = append(TeacherAvatarList, k)
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
	// 老师头像（file_id < 0 表示 -teacherId，走 teachers.avatar_url）
	// 批量收集所有需要查的 teacherId，一次查 teachers 表
	if len(TeacherAvatarList) > 0 {
		teacherIDs := make([]int, len(TeacherAvatarList))
		for i, tid := range TeacherAvatarList {
			teacherIDs[i] = -tid // 还原成正数 teacherId
		}
		type teacherAvatar struct {
			Id        int    `json:"id"`
			AvatarUrl string `json:"avatar_url"`
		}
		var avatars []teacherAvatar
		err = g.DB().Ctx(ctx).Model("teachers").
			WhereIn("id", teacherIDs).
			Fields("id, avatar_url").
			Scan(&avatars)
		if err != nil {
			return nil, err
		}
		teacherUrlMap := make(map[int]string, len(avatars))
		for _, a := range avatars {
			teacherUrlMap[a.Id] = a.AvatarUrl
		}
		for _, tid := range TeacherAvatarList {
			realId := -tid
			url := teacherUrlMap[realId]
			for _, aid := range AvatarsMap[tid] {
				answerList[IdMap[aid]].UserAvatar = url
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

	// 更新通知
	_, err = service.Notification().UpdateAoQ(ctx, model.UpdateAoQInput{UserID: userId, QuestionID: req.QuestionID})
	if err != nil {
		return nil, err
	}
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
