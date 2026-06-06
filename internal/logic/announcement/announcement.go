package announcement

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/custom"
	"suask/internal/service"
	"suask/utility"
	"suask/utility/files"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sAnnouncement struct{}

const pageSize = 10

func timestampMilli(t *gtime.Time) int64 {
	if t == nil {
		return 0
	}
	return t.TimestampMilli()
}

func countRemainPage(total, page, size int) int {
	remainNum := total - size*page
	if remainNum <= 0 {
		return 0
	}
	remain := remainNum / size
	if remainNum%size > 0 {
		remain++
	}
	return remain
}

func (s *sAnnouncement) List(ctx context.Context, in model.AnnouncementListInput) (*model.AnnouncementListOutput, error) {
	md := g.DB().Ctx(ctx).Model("announcements a").
		LeftJoin("users u", "u.id = a.author_id").
		Fields("a.id, a.title, a.contents, u.nickname AS author_name, a.is_pinned, a.published_at, a.expires_at").
		Where("a.deleted_at IS NULL").
		Where("a.expires_at IS NULL OR a.expires_at > ?", gtime.Now()).
		Order("a.is_pinned DESC, a.published_at DESC").
		Page(in.Page, pageSize)

	type row struct {
		Id          int         `json:"id"`
		Title       string      `json:"title"`
		Contents    string      `json:"contents"`
		AuthorName  string      `json:"author_name"`
		IsPinned    int         `json:"is_pinned"`
		PublishedAt *gtime.Time `json:"published_at"`
		ExpiresAt   *gtime.Time `json:"expires_at"`
	}
	var rows []row
	var total int
	if err := md.ScanAndCount(&rows, &total, false); err != nil {
		return nil, err
	}

	// 批量拿评论数
	ids := make([]int, len(rows))
	for i, r := range rows {
		ids[i] = r.Id
	}
	commentCounts := make(map[int]int)
	if len(ids) > 0 {
		type cnt struct {
			AnnouncementId int `json:"announcement_id"`
			Cnt            int `json:"cnt"`
		}
		var counts []cnt
		err := g.DB().Ctx(ctx).Model("answers").
			Fields("announcement_id, COUNT(*) AS cnt").
			WhereIn("announcement_id", ids).
			Where("deleted_at IS NULL").
			Group("announcement_id").
			Scan(&counts)
		if err == nil {
			for _, c := range counts {
				commentCounts[c.AnnouncementId] = c.Cnt
			}
		}
	}

	items := make([]model.AnnouncementListItem, len(rows))
	for i, r := range rows {
		items[i] = model.AnnouncementListItem{
			ID:          r.Id,
			Title:       r.Title,
			Content:     utility.TruncateString(r.Contents),
			AuthorName:  r.AuthorName,
			IsPinned:    r.IsPinned == 1,
			PublishedAt: timestampMilli(r.PublishedAt),
			ExpiresAt:   timestampMilli(r.ExpiresAt),
			CommentCnt:  commentCounts[r.Id],
		}
	}

	remain := countRemainPage(total, in.Page, pageSize)
	return &model.AnnouncementListOutput{Items: items, RemainPage: remain, Total: total}, nil
}

func (s *sAnnouncement) GetActive(ctx context.Context) (*model.AnnouncementActiveOutput, error) {
	type row struct {
		Id          int         `json:"id"`
		Title       string      `json:"title"`
		IsPinned    int         `json:"is_pinned"`
		PublishedAt *gtime.Time `json:"published_at"`
		ExpiresAt   *gtime.Time `json:"expires_at"`
	}

	var r row
	err := g.DB().Ctx(ctx).Model("announcements a").
		Fields("a.id, a.title, a.is_pinned, a.published_at, a.expires_at").
		Where("a.deleted_at IS NULL").
		Where("a.expires_at IS NULL OR a.expires_at > ?", gtime.Now()).
		Order("a.is_pinned DESC, a.published_at DESC").
		Limit(1).
		Scan(&r)
	if err != nil {
		return nil, err
	}
	if r.Id == 0 {
		return &model.AnnouncementActiveOutput{}, nil
	}
	return &model.AnnouncementActiveOutput{Item: &model.AnnouncementListItem{
		ID:          r.Id,
		Title:       r.Title,
		IsPinned:    r.IsPinned == 1,
		PublishedAt: timestampMilli(r.PublishedAt),
		ExpiresAt:   timestampMilli(r.ExpiresAt),
	}}, nil
}

func (s *sAnnouncement) Detail(ctx context.Context, in model.AnnouncementDetailInput) (*model.AnnouncementDetailOutput, error) {
	type row struct {
		Id          int         `json:"id"`
		Title       string      `json:"title"`
		Contents    string      `json:"contents"`
		AuthorName  string      `json:"author_name"`
		IsPinned    int         `json:"is_pinned"`
		PublishedAt *gtime.Time `json:"published_at"`
		ExpiresAt   *gtime.Time `json:"expires_at"`
	}
	var r row
	err := g.DB().Ctx(ctx).Model("announcements a").
		LeftJoin("users u", "u.id = a.author_id").
		Fields("a.id, a.title, a.contents, u.nickname AS author_name, a.is_pinned, a.published_at, a.expires_at").
		Where("a.id = ? AND a.deleted_at IS NULL", in.ID).
		Scan(&r)
	if err != nil {
		return nil, err
	}

	// 附件图片
	var imgs []custom.Image
	dao.Attachments.Ctx(ctx).Where("announcement_id = ?", in.ID).Scan(&imgs)
	imgIDs := make([]int, 0, len(imgs))
	for _, img := range imgs {
		imgIDs = append(imgIDs, img.FileID)
	}

	return &model.AnnouncementDetailOutput{
		ID:          r.Id,
		Title:       r.Title,
		Content:     r.Contents,
		AuthorName:  r.AuthorName,
		IsPinned:    r.IsPinned == 1,
		PublishedAt: timestampMilli(r.PublishedAt),
		ExpiresAt:   timestampMilli(r.ExpiresAt),
		ImageIDs:    imgIDs,
	}, nil
}

func (s *sAnnouncement) Create(ctx context.Context, in model.AnnouncementCreateInput) (*model.AnnouncementCreateOutput, error) {
	data := g.Map{
		"author_id":  in.AuthorID,
		"title":      in.Title,
		"contents":   in.Content,
		"is_pinned":  in.IsPinned,
		"expires_at": in.ExpiresAt,
	}
	id, err := g.DB().Ctx(ctx).Model("announcements").InsertAndGetId(data)
	if err != nil {
		return nil, err
	}
	return &model.AnnouncementCreateOutput{ID: int(id)}, nil
}

func (s *sAnnouncement) Update(ctx context.Context, in model.AnnouncementUpdateInput) (*model.AnnouncementUpdateOutput, error) {
	data := g.Map{"updated_at": gtime.Now()}
	if in.Title != "" {
		data["title"] = in.Title
	}
	if in.Content != "" {
		data["contents"] = in.Content
	}
	if in.IsPinned != nil {
		data["is_pinned"] = *in.IsPinned
	}
	if in.ClearExpiresAt {
		data["expires_at"] = gdb.Raw("NULL")
	} else if in.ExpiresAt != nil {
		data["expires_at"] = in.ExpiresAt
	}
	_, err := g.DB().Ctx(ctx).Model("announcements").Where("id = ?", in.ID).Update(data)
	if err != nil {
		return nil, err
	}
	return &model.AnnouncementUpdateOutput{ID: in.ID}, nil
}

func (s *sAnnouncement) Delete(ctx context.Context, in model.AnnouncementDeleteInput) error {
	_, err := g.DB().Ctx(ctx).Model("announcements").Where("id = ?", in.ID).Update(g.Map{"deleted_at": gtime.Now()})
	return err
}

func (s *sAnnouncement) AddComment(ctx context.Context, in model.AnnouncementCommentInput) (*model.AnnouncementCommentOutput, error) {
	data := g.Map{
		"user_id":         in.UserID,
		"announcement_id": in.AnnouncementID,
		"contents":        in.Content,
		"in_reply_to":     in.InReplyTo,
	}
	id, err := g.DB().Ctx(ctx).Model("answers").InsertAndGetId(data)
	if err != nil {
		return nil, err
	}
	return &model.AnnouncementCommentOutput{ID: int(id)}, nil
}

func (s *sAnnouncement) GetComments(ctx context.Context, announcementID int) ([]model.AnnouncementComment, error) {
	type row struct {
		Id        int         `json:"id"`
		UserId    int         `json:"user_id"`
		Nickname  string      `json:"nickname"`
		AvatarId  int         `json:"avatar_file_id"`
		Contents  string      `json:"contents"`
		CreatedAt *gtime.Time `json:"created_at"`
		InReplyTo int         `json:"in_reply_to"`
	}
	var rows []row
	err := g.DB().Ctx(ctx).Model("answers a").
		LeftJoin("users u", "u.id = a.user_id").
		Fields("a.id, a.user_id, u.nickname, u.avatar_file_id, a.contents, a.created_at, a.in_reply_to").
		Where("a.announcement_id = ? AND a.deleted_at IS NULL", announcementID).
		Order("a.created_at ASC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	// 批量查头像 URL（#5 修复：N+1 → 1）
	avatarIDs := make([]int, 0, len(rows))
	for _, r := range rows {
		if r.AvatarId != 0 {
			avatarIDs = append(avatarIDs, r.AvatarId)
		}
	}
	urlMap := make(map[int]string)
	if len(avatarIDs) > 0 {
		var fileList []struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
			Hash []byte `json:"hash"`
		}
		dao.Files.Ctx(ctx).WhereIn("id", avatarIDs).Scan(&fileList)
		for _, f := range fileList {
			if url, err := files.GetURL(f.Hash, f.Name); err == nil {
				urlMap[f.Id] = url
			}
		}
	}

	comments := make([]model.AnnouncementComment, len(rows))
	for i, r := range rows {
		avatar := consts.DefaultAvatarURL
		if r.AvatarId != 0 {
			if url, ok := urlMap[r.AvatarId]; ok {
				avatar = url
			}
		}
		comments[i] = model.AnnouncementComment{
			ID:        r.Id,
			UserID:    r.UserId,
			Nickname:  r.Nickname,
			AvatarURL: avatar,
			Contents:  r.Contents,
			CreatedAt: r.CreatedAt.TimestampMilli(),
			InReplyTo: r.InReplyTo,
		}
	}
	return comments, nil
}

func init() {
	service.RegisterAnnouncement(New())
}

func New() *sAnnouncement {
	return &sAnnouncement{}
}
