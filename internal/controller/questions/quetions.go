package questions

import (
	"context"
	"fmt"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"
	"suask/utility/validation"

	"github.com/gogf/gf/v2/util/gconv"
)

type cQuestion struct{}

var Question = cQuestion{}

func (cQuestion) Add(ctx context.Context, req *v1.AddQuestionReq) (res *v1.AddQuestionRes, err error) {
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	//UserId := 2
	//if UserId == consts.DefaultUserId {
	//	return nil, fmt.Errorf("user not login")
	//}
	if req.DstUserId == 0 && UserId == consts.DefaultUserId {
		return nil, fmt.Errorf("未登录用户不能提问大家")
	}
	if req.DstUserId != 0 {
		// 防止非法提问
		err = validation.TeacherPerm(ctx, req.DstUserId)
		if err != nil {
			return
		}
	}
	questionInput := model.AddQuestionInput{}
	err = gconv.Struct(req, &questionInput)
	if err != nil {
		return nil, err
	}
	questionInput.SrcUserID = UserId
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
