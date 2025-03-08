// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"suask/internal/model"
)

type (
	IQuestionDetail interface {
		GetQuestionBase(ctx context.Context, in *model.GetQuestionBaseInput) (*model.GetQuestionBaseOutput, error)
		GetAnswers(ctx context.Context, in *model.GetAnswerDetailInput) (*model.GetAnswerDetailOutput, error)
		AddQuestionView(ctx context.Context, in *model.AddViewInput) (*model.AddViewOutput, error)
		AddAnswerUpvote(ctx context.Context, in *model.UpvoteInput) (*model.UpvoteOutput, error)
		ReplyQuestion(ctx context.Context, in *model.AddAnswerInput) (*model.AddAnswerOutput, error)
		AddReplyCnt(ctx context.Context, in *model.AddReplyCntInput) (*model.AddReplyCntOutput, error)
		BuildRelation(ctx context.Context, in *model.BuildRelationInput) (*model.BuildRelationOutput, error)
	}
)

var (
	localQuestionDetail IQuestionDetail
)

func QuestionDetail() IQuestionDetail {
	if localQuestionDetail == nil {
		panic("implement not found for interface IQuestionDetail, forgot register?")
	}
	return localQuestionDetail
}

func RegisterQuestionDetail(i IQuestionDetail) {
	localQuestionDetail = i
}
