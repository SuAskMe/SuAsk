package questions_detail

import (
	"context"
	"fmt"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/model/entity"
	"suask/internal/service"
)

type sQuestionDetail struct{}

func (sQuestionDetail) GetQuestionBase(ctx context.Context, in *model.GetQuestionBaseInput) (*model.GetQuestionBaseOutput, error) {
	UserId := 1
	// UserId := gconv.Int(ctx.Value(consts.CtxId))
	md := dao.Questions.Ctx(ctx).Where("id =?", in.QuestionId)
	var question entity.Questions
	err := md.Scan(&question)
	if err != nil {
		return nil, err
	}
	if question.IsPrivate && question.SrcUserId != UserId {
		return nil, fmt.Errorf("you are not allowed to access this question")
	}
	canReply := true
	if question.DstUserId != 0 && (question.DstUserId != UserId || question.SrcUserId != UserId) {
		// 如果目标用户不为空
		// 不是问自己，或者不是自己的提问，不能回复
		canReply = false
	}
	output := model.GetQuestionBaseOutput{
		Question: &model.QuestionBase{
			ID:        question.Id,
			Title:     question.Title,
			Content:   question.Contents,
			Views:     question.Views,
			CreatedAt: question.CreatedAt.TimestampMilli(),
		},
		CanReply: canReply,
	}
	return &output, nil
}

func (sQuestionDetail) GetAnswers(ctx context.Context, in *model.GetAnswerDetailInput) (*model.GetAnswerDetailOutput, error) {
	md := dao.Answers.Ctx(ctx).Where("question_id =?", in.QuestionId)
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
		if _, ok := UserIdMap[ans.UserId]; !ok {
			UserIdMap[ans.UserId] = []int{ans.Id}
		} else {
			UserIdMap[ans.UserId] = append(UserIdMap[ans.UserId], ans.Id)
		}
		IdMap[ans.Id] = i
		answerList[i] = model.AnswerWithDetails{
			Id:        ans.Id,
			UserId:    ans.UserId,
			Contents:  ans.Contents,
			CreatedAt: ans.CreatedAt.TimestampMilli(),
			Upvotes:   ans.Upvotes,
		}
	}
	// 获取回答者的头像
	UserIdList := make([]int, 0, len(answers)) // 用户ID列表
	for k := range UserIdMap {
		UserIdList = append(UserIdList, k)
	}
	md = dao.Users.Ctx(ctx).WhereIn("id IN (?)", UserIdList)
	var avatars []custom.Avatars // 用户头像
	err = md.Scan(&avatars)
	if err != nil {
		return nil, err
	}
	AvatarMap := make(map[int][]int) // 头像ID对应的回答ID列表
	for _, avatar := range avatars {
		AvatarMap[avatar.AvatarFileId] = UserIdMap[avatar.UserId]
	}
	// 获取回答的图片
	md = dao.Attachments.Ctx(ctx).WhereIn("answer_id IN (?)", IdList)
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
		Answers:    answerList,
		AvatarsMap: AvatarMap,
		ImageMap:   ImgMap,
	}, nil
}

func init() {
	service.RegisterQuestionDetail(New())
}

func New() *sQuestionDetail {
	return &sQuestionDetail{}
}
