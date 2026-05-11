package questions

import (
	"context"
	"fmt"
	v1 "suask/api/questions/v1"
	"suask/internal/consts"
	qutil "suask/internal/logic/questions_util"
	"suask/internal/model"
	"suask/internal/service"
	"suask/module/validation"

	"github.com/gogf/gf/v2/util/gconv"
)

// cInbox 统一收件箱 controller，替代旧的 /teacher/question/* 系列。
type cInbox struct{}

var Inbox = cInbox{}

func (cInbox) Get(ctx context.Context, req *v1.InboxReq) (res *v1.InboxRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if uid == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录")
	}
	_, err = validation.IsTeacher(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("仅老师可访问收件箱")
	}

	tag := req.Tag
	if tag == "" {
		tag = "all"
	}
	in := model.GetQFMInput{
		TeacherId: uid,
		Page:      req.Page,
		SortType:  req.SortType,
	}
	switch tag {
	case "answered":
		in.Tag = consts.Answered
	case "unanswered":
		in.Tag = consts.Unanswered
	case "pinned":
		out, err := service.TeacherQuestionSelf().GetQFMPinned(ctx, &in)
		if err != nil {
			return nil, err
		}
		return buildInboxRes(ctx, out)
	}

	out, err := service.TeacherQuestionSelf().GetQFMAll(ctx, &in)
	if err != nil {
		return nil, err
	}
	return buildInboxRes(ctx, out)
}

func (cInbox) Keywords(ctx context.Context, req *v1.InboxKeywordsReq) (res *v1.InboxKeywordsRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if uid == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录")
	}
	_, err = validation.IsTeacher(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("仅老师可访问")
	}
	in := model.GetQFMKeywordsInput{
		TeacherId: uid,
		Keyword:   req.Keyword,
	}
	out, err := service.TeacherQuestionSelf().GetKeyword(ctx, &in)
	if err != nil {
		return nil, err
	}
	gconv.Scan(out, &res)
	return
}

func (cInbox) Search(ctx context.Context, req *v1.InboxSearchReq) (res *v1.InboxSearchRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if uid == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录")
	}
	_, err = validation.IsTeacher(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("仅老师可访问")
	}
	in := model.GetQFMInput{
		TeacherId: uid,
		Page:      req.Page,
		SortType:  req.SortType,
		Keyword:   req.Keyword,
	}
	out, err := service.TeacherQuestionSelf().GetQFMAll(ctx, &in)
	if err != nil {
		return nil, err
	}
	inboxRes, err := buildInboxRes(ctx, out)
	if err != nil {
		return nil, err
	}
	return &v1.InboxSearchRes{
		QuestionList: inboxRes.QuestionList,
		RemainPage:   inboxRes.RemainPage,
	}, nil
}

func (cInbox) Pin(ctx context.Context, req *v1.PinReq) (res *v1.PinRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if uid == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录")
	}
	_, err = validation.IsTeacher(ctx, uid)
	if err != nil {
		return nil, fmt.Errorf("仅老师可操作")
	}
	out, err := service.TeacherQuestionSelf().PinQFM(ctx, &model.PinQFMInput{
		QuestionId: req.QuestionId,
		TeacherId:  uid,
	})
	if err != nil {
		return nil, err
	}
	return &v1.PinRes{IsPinned: out.IsPinned}, nil
}

// buildInboxRes 把 logic 层的 QFMOutput 转成 API 响应，顺带批量查图片。
func buildInboxRes(ctx context.Context, out *model.GetQFMOutput) (*v1.InboxRes, error) {
	qfm := out.Questions
	idMap := out.IdMap
	imagesOutput, err := service.QuestionUtil().GetImages(ctx, &model.GetImagesInput{QuestionIDs: out.QuestionIDs})
	if err != nil {
		return nil, err
	}
	if imagesOutput != nil && len(imagesOutput.ImageMap) > 0 {
		allFileIDs := qutil.CollectFileIDs(imagesOutput.ImageMap, nil)
		urlMap, err := qutil.BatchGetFileURLs(ctx, allFileIDs)
		if err != nil {
			return nil, err
		}
		imageURLs := qutil.ResolveImageURLs(urlMap, imagesOutput.ImageMap)
		for qid, urls := range imageURLs {
			qfm[idMap[qid]].ImageURLs = urls
		}
	}
	return &v1.InboxRes{
		QuestionList: qfm,
		RemainPage:   out.RemainPage,
	}, nil
}
