package attachements

import (
	"context"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"
)

type sAttachment struct{}

func (s *sAttachment) AddAttachments(ctx context.Context, in model.AddAttachmentInput) (out model.AddAttachmentOutput, err error) {
	fileCount := len(in.FileId)
	attachment := do.Attachments{
		QuestionId: in.QuestionId,
		AnswerId:   in.AnswerId,
		Type:       in.Type,
	}
	out = model.AddAttachmentOutput{
		Id: make([]int, fileCount),
	}
	for i := 0; i < fileCount; i++ {
		attachment.FileId = in.FileId[i]
		id, err := dao.Attachments.Ctx(ctx).InsertAndGetId(attachment)
		if err != nil {
			return model.AddAttachmentOutput{}, err
		}
		out.Id[i] = int(id)
	}
	return out, nil
}

func init() {
	service.RegisterAttachment(New())
}

func New() *sAttachment {
	return &sAttachment{}
}
