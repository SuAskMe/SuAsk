package questions

import (
	"context"
	"fmt"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/send_email"
	"suask/module/validation"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type cQuestion struct{}

var Question = cQuestion{}

func (cQuestion) Add(ctx context.Context, req *v1.AddQuestionReq) (res *v1.AddQuestionRes, err error) {
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	// "问大家"已下线：dst_user_id 必填
	if req.DstUserId == 0 {
		return nil, fmt.Errorf("必须指定目标老师")
	}
	// 防止非法提问（不存在的老师 / 权限不足）
	if err = validation.TeacherPerm(ctx, req.DstUserId); err != nil {
		return nil, err
	}
	questionInput := model.AddQuestionInput{}
	err = gconv.Struct(req, &questionInput)
	if err != nil {
		return nil, err
	}
	questionInput.SrcUserID = UserId
	var fileIdList model.FileListAddOutput
	if req.Files != nil {
		fileList := model.FileListAddInput{
			UploaderId: UserId,
			FileList:   req.Files,
		}
		fileIdList, err = service.File().UploadFileList(ctx, fileList)
		if err != nil {
			return nil, err
		}
	}

	var questionOut *model.AddQuestionOutput
	err = g.DB().Transaction(ctx, func(ctx context.Context, _ gdb.TX) error {
		questionOut, err = service.QuestionUtil().AddQuestion(ctx, &questionInput)
		if err != nil {
			return err
		}
		_, err = service.Attachment().AddAttachments(ctx, model.AddAttachmentInput{
			QuestionId: questionOut.ID,
			Type:       consts.QuestionFileType,
			FileId:     fileIdList.IdList,
		})
		if err != nil {
			return err
		}
		_, err = service.Notification().Add(ctx, model.AddNotificationInput{
			UserId:     req.DstUserId,
			QuestionId: questionOut.ID,
			Type:       consts.NewQuestion,
		})
		return err
	})
	if err != nil {
		return nil, err
	}
	res = &v1.AddQuestionRes{Id: questionOut.ID}

	err = service.Notification().SendNoticeEmail(ctx, &model.SendNoticeEmailInput{
		To: req.DstUserId,
		Notice: &send_email.Notice{
			User:    "SuAsk用户",
			Type:    "新的提问",
			Content: req.Content,
			URL:     "https://suask.me/question-detail/" + gconv.String(questionOut.ID),
		},
	})
	if err != nil {
		g.Log().Warningf(ctx, "发送提问邮件通知失败: %v", err)
	}
	return res, nil
}

func (cQuestion) Delete(ctx context.Context, req *v1.DeleteQuestionReq) (res *v1.DeleteQuestionRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录后操作")
	}
	if err := service.QuestionDetail().DeleteQuestion(ctx, req.ID, userId); err != nil {
		return nil, err
	}
	return &v1.DeleteQuestionRes{}, nil
}

func (cQuestion) Restore(ctx context.Context, req *v1.RestoreQuestionReq) (res *v1.RestoreQuestionRes, err error) {
	userId := gconv.Int(ctx.Value(consts.CtxId))
	if userId == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录后操作")
	}
	if err := service.QuestionDetail().RestoreQuestion(ctx, req.ID, userId); err != nil {
		return nil, err
	}
	return &v1.RestoreQuestionRes{}, nil
}
