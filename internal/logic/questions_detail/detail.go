package questions_detail

import (
	"context"
	"fmt"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
	"suask/module/validation"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sQuestionDetail struct{}

func AddResponseCnt(ctx context.Context, teacherId int) error {
	md := dao.Teachers.Ctx(ctx).Where(dao.Teachers.Columns().Id, teacherId)
	_, err := md.Increment(dao.Teachers.Columns().Responses, 1)
	if err != nil {
		return err
	}
	return nil
}

func (sQuestionDetail) GetQuestionBase(ctx context.Context, in *model.GetQuestionBaseInput) (*model.GetQuestionBaseOutput, error) {
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().Id, in.QuestionId)
	var question entity.Questions
	err := md.Scan(&question)
	if err != nil {
		return nil, err
	}
	// 权限验证
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	if question.DstUserId != 0 && question.DstUserId != UserId { // 问老师的问题, 且不是老师
		err = validation.TeacherPerm(ctx, question.DstUserId)
		if err != nil {
			return nil, err
		}
	}
	err = validation.QuestionPerm(ctx, &question)
	if err != nil {
		return nil, err
	}
	canReply := true
	err = validation.AnswerPerm(ctx, &question)
	if err != nil {
		canReply = false
	}
	// 获取问题详情
	var imgList []custom.Image
	var count int
	err = dao.Attachments.Ctx(ctx).Where(dao.Attachments.Columns().QuestionId, question.Id).ScanAndCount(&imgList, &count, false)
	if err != nil {
		return nil, err
	}
	imgIdList := make([]int, len(imgList))
	for _, img := range imgList {
		imgIdList = append(imgIdList, img.FileID)
	}
	isFavorite := false
	if in.UserId != consts.DefaultUserId {
		one, err := dao.Favorites.Ctx(ctx).Where(dao.Favorites.Columns().QuestionId, in.QuestionId).Where(dao.Favorites.Columns().UserId, in.UserId).One()
		if !one.IsEmpty() {
			isFavorite = true
		}
		if err != nil {
			return nil, err
		}
	}
	output := model.GetQuestionBaseOutput{
		ID:         question.Id,
		Title:      question.Title,
		Content:    question.Contents,
		Views:      question.Views,
		CreatedAt:  question.CreatedAt.TimestampMilli(),
		CanReply:   canReply,
		ImageList:  imgIdList,
		IsFavorite: isFavorite,
		DstUserId:  question.DstUserId,
	}
	return &output, nil
}

func (sQuestionDetail) GetAnswers(ctx context.Context, in *model.GetAnswerDetailInput) (*model.GetAnswerDetailOutput, error) {
	md := dao.Answers.Ctx(ctx).Where(dao.Answers.Columns().QuestionId, in.QuestionId)
	var answers []entity.Answers
	err := md.Scan(&answers)
	if err != nil {
		return nil, err
	}
	// 获取回答详情

	answerList := make([]model.AnswerWithDetails, len(answers))
	IdList := make([]int, len(answers)) // 回答的ID列表
	IdMap := make(map[int]int)          // 回答的ID映射
	UserIdMap := make(map[int][]int)    // 用户ID所对应的回答ID列表

	for i, ans := range answers {
		IdList[i] = ans.Id
		IdMap[ans.Id] = i
		UserId := ans.UserId

		if in.DstUserId != 0 && in.DstUserId != UserId { // 问老师的问题且不是老师的回答
			UserId = consts.DefaultUserId // 显示默认用户的头像
		}

		if _, ok := UserIdMap[ans.UserId]; !ok {
			UserIdMap[UserId] = []int{ans.Id}
		} else {
			UserIdMap[UserId] = append(UserIdMap[UserId], ans.Id)
		}

		answerList[i].Id = ans.Id
		answerList[i].UserId = UserId
		answerList[i].InReplyTo = ans.InReplyTo
		answerList[i].Contents = ans.Contents
		answerList[i].CreatedAt = ans.CreatedAt.TimestampMilli()
		answerList[i].Upvotes = ans.Upvotes
	}

	// 获取用户点赞信息
	//UserId := 2
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	md = dao.Upvotes.Ctx(ctx).WhereIn(dao.Upvotes.Columns().AnswerId, IdList).Where(dao.Upvotes.Columns().UserId, UserId)
	var upvotes []custom.MyUpvotes
	err = md.Scan(&upvotes)
	if err != nil {
		return nil, err
	}
	for _, upvote := range upvotes {
		answerList[IdMap[upvote.AnswerId]].IsUpvoted = true
	}
	// 获取回答者的信息
	UserIdList := make([]int, 0, len(answers)) // 用户ID列表
	for k := range UserIdMap {
		UserIdList = append(UserIdList, k)
	}
	md = dao.Users.Ctx(ctx).WhereIn(dao.Users.Columns().Id, UserIdList)
	var userInfo []custom.UserInfo // 用户信息
	err = md.Scan(&userInfo)
	if err != nil {
		return nil, err
	}
	AvatarMap := make(map[int][]int) // 头像ID对应的回答ID列表
	for _, info := range userInfo {
		for _, v := range UserIdMap[info.UserId] {
			answerList[IdMap[v]].NickName = info.NickName
		}
		if info.Role == consts.TEACHER { // 如果是老师，则显示用户名
			for _, v := range UserIdMap[info.UserId] {
				answerList[IdMap[v]].TeacherName = info.Name
			}
		}
		if info.AvatarFileId == 0 && info.Role == consts.TEACHER {
			// 老师没有头像，显示信息头像
			AvatarMap[-info.UserId] = UserIdMap[info.UserId]
		} else {
			AvatarMap[info.AvatarFileId] = UserIdMap[info.UserId]
		}

	}
	// 获取回答的图片
	md = dao.Attachments.Ctx(ctx).WhereIn(dao.Attachments.Columns().AnswerId, IdList)
	var imgList []custom.AnswerImage
	err = md.Scan(&imgList)
	if err != nil {
		return nil, err
	}
	ImgMap := make(map[int][]int)
	for _, img := range imgList {
		if _, ok := ImgMap[img.AnswerId]; !ok {
			ImgMap[img.AnswerId] = make([]int, 0, 8)
		}
		ImgMap[img.AnswerId] = append(ImgMap[img.AnswerId], img.FileID)
	}
	return &model.GetAnswerDetailOutput{
		IdMap:      IdMap,
		Answers:    answerList,
		AvatarsMap: AvatarMap,
		ImageMap:   ImgMap,
	}, nil
}

func (sQuestionDetail) AddQuestionView(ctx context.Context, in *model.AddViewInput) (*model.AddViewOutput, error) {
	db := g.DB()
	// ViewsLock.Lock()
	_, err := db.Exec(ctx, "UPDATE questions SET views = views + 1 WHERE id = ?", in.QuestionId)
	// ViewsLock.Unlock()
	if err != nil {
		return nil, err
	}
	return &model.AddViewOutput{}, nil
}

func (sQuestionDetail) AddAnswerUpvote(ctx context.Context, in *model.UpvoteInput) (*model.UpvoteOutput, error) {
	//UserId := 2
	db := g.DB()
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	if UserId == consts.DefaultUserId {
		return nil, fmt.Errorf("you are not allowed to upvote")
	}
	md := dao.Upvotes.Ctx(ctx).Where(do.Upvotes{AnswerId: in.AnswerId, UserId: UserId})
	cnt, err := md.Count()
	if err != nil {
		return nil, err
	}
	if cnt > 0 {
		_, err = md.Delete()
		if err != nil {
			return nil, err
		}
		md = dao.Answers.Ctx(ctx).Where(dao.Answers.Columns().Id, in.AnswerId)
		_, err = db.Exec(ctx, "UPDATE answers SET upvotes = upvotes - 1 WHERE id = ?", in.AnswerId)
		if err != nil {
			return nil, err
		}
		res, err := md.One()
		if err != nil {
			return nil, err
		}
		cnt = res["upvotes"].Int()
		return &model.UpvoteOutput{
			IsUpvoted: false,
			Upvotes:   cnt,
		}, nil
	} else {
		md = dao.Upvotes.Ctx(ctx)
		_, err = md.Insert(do.Upvotes{AnswerId: in.AnswerId, UserId: UserId})
		if err != nil {
			return nil, err
		}
		md = dao.Answers.Ctx(ctx).Where("id=?", in.AnswerId)
		_, err = db.Exec(ctx, "UPDATE answers SET upvotes = upvotes + 1 WHERE id = ?", in.AnswerId)
		if err != nil {
			return nil, err
		}
		res, err := md.One()
		if err != nil {
			return nil, err
		}
		cnt = res["upvotes"].Int()
		return &model.UpvoteOutput{
			IsUpvoted: true,
			Upvotes:   cnt,
		}, nil
	}
}

func (sQuestionDetail) ReplyQuestion(ctx context.Context, in *model.AddAnswerInput) (*model.AddAnswerOutput, error) {
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().Id, in.QuestionId).Fields("id, is_private, src_user_id, dst_user_id, reply_cnt")
	var question entity.Questions
	err := md.Scan(&question)
	if err != nil {
		return nil, err
	}
	// 权限验证
	if question.DstUserId != 0 && question.DstUserId != in.UserId { // 问老师的问题且不是老师的回答
		err = validation.TeacherPerm(ctx, question.DstUserId)
		if err != nil {
			return nil, err
		}
	}
	err = validation.AnswerPerm(ctx, &question)
	if err != nil {
		return nil, fmt.Errorf("you are not allowed to access this question")
	}
	// 保存回答
	md = dao.Answers.Ctx(ctx)
	id, err := md.InsertAndGetId(do.Answers{
		QuestionId: in.QuestionId,
		UserId:     in.UserId,
		Contents:   in.Content,
		InReplyTo:  in.InReplyTo,
	})
	if err != nil {
		return nil, err
	}

	if question.DstUserId != 0 && question.ReplyCnt == 0 { // 问老师的问题，且第一次回复，增加回复数
		err = AddResponseCnt(ctx, question.DstUserId)
		if err != nil {
			return nil, err
		}
	}

	return &model.AddAnswerOutput{
		Id: int(id),
	}, nil
}

func (sQuestionDetail) AddReplyCnt(ctx context.Context, in *model.AddReplyCntInput) (*model.AddReplyCntOutput, error) {
	md := dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().Id, in.QuestionId)
	_, err := md.Increment(dao.Questions.Columns().ReplyCnt, 1)
	if err != nil {
		return nil, err
	}
	md = dao.Questions.Ctx(ctx).Where(dao.Questions.Columns().Id, in.QuestionId).Fields(dao.Questions.Columns().ReplyCnt)
	res, err := md.One()
	if err != nil {
		return nil, err
	}
	cnt := res[dao.Questions.Columns().ReplyCnt].Int()
	return &model.AddReplyCntOutput{ReplyCnt: cnt}, nil
}

func (sQuestionDetail) BuildRelation(ctx context.Context, in *model.BuildRelationInput) (*model.BuildRelationOutput, error) {
	// 保存关系
	UserId := gconv.Int(ctx.Value(consts.CtxId))
	md := dao.UserRelation.Ctx(ctx)
	_, err := md.Insert(do.UserRelation{
		QuestionId: in.QuestionId,
		UserId:     UserId,
	})
	if err != nil {
		return nil, err
	}
	return &model.BuildRelationOutput{}, nil
}

func init() {
	service.RegisterQuestionDetail(New())
}

func New() *sQuestionDetail {
	return &sQuestionDetail{}
}
