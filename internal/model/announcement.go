package model

import "github.com/gogf/gf/v2/os/gtime"

// --- Service 层入参/出参 ---

type AnnouncementListInput struct {
	Page           int
	IncludeExpired bool
}

type AnnouncementListOutput struct {
	Items      []AnnouncementListItem
	RemainPage int
	Total      int
}

type AnnouncementListItem struct {
	ID          int
	Title       string
	Content     string
	AuthorName  string
	IsPinned    bool
	PublishedAt int64
	ExpiresAt   int64
	CommentCnt  int
}

type AnnouncementDetailInput struct {
	ID int
}

type AnnouncementDetailOutput struct {
	ID          int
	Title       string
	Content     string
	AuthorName  string
	IsPinned    bool
	PublishedAt int64
	ExpiresAt   int64
	ImageIDs    []int
}

type AnnouncementCreateInput struct {
	AuthorID  int
	Title     string
	Content   string
	IsPinned  bool
	ExpiresAt *gtime.Time
}

type AnnouncementCreateOutput struct {
	ID int
}

type AnnouncementUpdateInput struct {
	ID             int
	Title          string
	Content        string
	IsPinned       *bool
	ExpiresAt      *gtime.Time
	ClearExpiresAt bool
}

type AnnouncementUpdateOutput struct {
	ID int
}

type AnnouncementDeleteInput struct {
	ID int
}

type AnnouncementCommentInput struct {
	AnnouncementID int
	UserID         int
	Content        string
	InReplyTo      *int
}

type AnnouncementCommentOutput struct {
	ID int
}

type AnnouncementComment struct {
	ID        int
	UserID    int
	Nickname  string
	AvatarURL string
	Contents  string
	CreatedAt int64
	InReplyTo int
}

type AnnouncementActiveOutput struct {
	Item *AnnouncementListItem
}
