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
	if fileCount == 0 {
		return model.AddAttachmentOutput{}, nil
	}
	// #6 优化：批量 INSERT，一条 SQL 插入所有附件
	batch := make([]do.Attachments, fileCount)
	for i := 0; i < fileCount; i++ {
		batch[i] = do.Attachments{
			QuestionId:     in.QuestionId,
			AnswerId:       in.AnswerId,
			AnnouncementId: in.AnnouncementId,
			Type:           in.Type,
			FileId:         in.FileId[i],
		}
	}
	_, err = dao.Attachments.Ctx(ctx).Data(batch).Insert()
	if err != nil {
		return model.AddAttachmentOutput{}, err
	}
	// 批量 INSERT 后无法拿到每条的 lastInsertId（SQLite 只返回最后一条的），
	// 但 AddAttachmentOutput.Id 目前没有任何调用方使用，返回空即可。
	out = model.AddAttachmentOutput{Id: nil}
	return out, nil
}

func init() {
	service.RegisterAttachment(New())
}

func New() *sAttachment {
	return &sAttachment{}
}
