package service

import (
	"context"
	"suask/internal/model"
)

type (
	IAnnouncement interface {
		List(ctx context.Context, in model.AnnouncementListInput) (*model.AnnouncementListOutput, error)
		Detail(ctx context.Context, in model.AnnouncementDetailInput) (*model.AnnouncementDetailOutput, error)
		GetActive(ctx context.Context) (*model.AnnouncementActiveOutput, error)
		Create(ctx context.Context, in model.AnnouncementCreateInput) (*model.AnnouncementCreateOutput, error)
		Update(ctx context.Context, in model.AnnouncementUpdateInput) (*model.AnnouncementUpdateOutput, error)
		SyncImages(ctx context.Context, in model.AnnouncementImageSyncInput) error
		Delete(ctx context.Context, in model.AnnouncementDeleteInput) error
		AddComment(ctx context.Context, in model.AnnouncementCommentInput) (*model.AnnouncementCommentOutput, error)
		GetComments(ctx context.Context, announcementID int) ([]model.AnnouncementComment, error)
	}
)

var localAnnouncement IAnnouncement

func Announcement() IAnnouncement {
	if localAnnouncement == nil {
		panic("implement not found for interface IAnnouncement, forgot register?")
	}
	return localAnnouncement
}

func RegisterAnnouncement(i IAnnouncement) {
	localAnnouncement = i
}
