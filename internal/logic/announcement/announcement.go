package announcement

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
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

type announcementImageAttachment struct {
	Id             int `json:"id"`
	AnnouncementId int `json:"announcement_id"`
	FileId         int `json:"file_id"`
}

func (s *sAnnouncement) getImageURLs(ctx context.Context, announcementIDs []int) (map[int][]string, map[int][]model.AnnouncementImage, error) {
	imageURLs := make(map[int][]string, len(announcementIDs))
	images := make(map[int][]model.AnnouncementImage, len(announcementIDs))
	if len(announcementIDs) == 0 {
		return imageURLs, images, nil
	}

	var attachments []announcementImageAttachment
	err := g.DB().Ctx(ctx).Model("attachments").
		Fields("id, announcement_id, file_id").
		WhereIn("announcement_id", announcementIDs).
		Where("type = ?", consts.QuestionFileType).
		Order("id ASC").
		Scan(&attachments)
	if err != nil {
		return nil, nil, err
	}
	if len(attachments) == 0 {
		return imageURLs, images, nil
	}

	fileIDs := make([]int, 0, len(attachments))
	seen := make(map[int]struct{}, len(attachments))
	for _, attachment := range attachments {
		if attachment.FileId == 0 {
			continue
		}
		if _, ok := seen[attachment.FileId]; ok {
			continue
		}
		seen[attachment.FileId] = struct{}{}
		fileIDs = append(fileIDs, attachment.FileId)
	}
	fileList, err := service.File().GetList(ctx, model.FileListGetInput{IdList: fileIDs})
	if err != nil {
		return nil, nil, err
	}
	urlByFileID := make(map[int]string, len(fileList.FileId))
	for i, fileID := range fileList.FileId {
		if i < len(fileList.URL) {
			urlByFileID[fileID] = fileList.URL[i]
		}
	}

	for _, attachment := range attachments {
		url, ok := urlByFileID[attachment.FileId]
		if !ok || url == "" {
			continue
		}
		imageURLs[attachment.AnnouncementId] = append(imageURLs[attachment.AnnouncementId], url)
		images[attachment.AnnouncementId] = append(images[attachment.AnnouncementId], model.AnnouncementImage{
			ID:  attachment.FileId,
			URL: url,
		})
	}
	return imageURLs, images, nil
}

func (s *sAnnouncement) List(ctx context.Context, in model.AnnouncementListInput) (*model.AnnouncementListOutput, error) {
	md := g.DB().Ctx(ctx).Model("announcements a").
		LeftJoin("users u", "u.id = a.author_id").
		Fields("a.id, a.title, a.contents, u.nickname AS author_name, a.is_pinned, a.published_at, a.expires_at").
		Where("a.deleted_at IS NULL").
		Order("a.is_pinned DESC, a.published_at DESC").
		Page(in.Page, pageSize)
	if !in.IncludeExpired {
		md = md.Where("a.expires_at IS NULL OR a.expires_at > ?", gtime.Now())
	}

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

	imageURLs, _, err := s.getImageURLs(ctx, ids)
	if err != nil {
		return nil, err
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
			ImageURLs:   imageURLs[r.Id],
		}
	}

	remain := countRemainPage(total, in.Page, pageSize)
	return &model.AnnouncementListOutput{Items: items, RemainPage: remain, Total: total}, nil
}

func (s *sAnnouncement) GetActive(ctx context.Context) (*model.AnnouncementActiveOutput, error) {
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
		Where("a.deleted_at IS NULL").
		Where("a.is_pinned = ?", 1).
		Where("a.expires_at IS NULL OR a.expires_at > ?", gtime.Now()).
		Order("a.published_at DESC").
		Limit(1).
		Scan(&r)
	if err != nil {
		return nil, err
	}
	if r.Id == 0 {
		return &model.AnnouncementActiveOutput{}, nil
	}
	imageURLs, _, err := s.getImageURLs(ctx, []int{r.Id})
	if err != nil {
		return nil, err
	}
	return &model.AnnouncementActiveOutput{Item: &model.AnnouncementActiveItem{
		ID:          r.Id,
		Title:       r.Title,
		Content:     r.Contents,
		AuthorName:  r.AuthorName,
		IsPinned:    r.IsPinned == 1,
		PublishedAt: timestampMilli(r.PublishedAt),
		ExpiresAt:   timestampMilli(r.ExpiresAt),
		ImageURLs:   imageURLs[r.Id],
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

	imageURLs, images, err := s.getImageURLs(ctx, []int{in.ID})
	if err != nil {
		return nil, err
	}

	return &model.AnnouncementDetailOutput{
		ID:          r.Id,
		Title:       r.Title,
		Content:     r.Contents,
		AuthorName:  r.AuthorName,
		IsPinned:    r.IsPinned == 1,
		PublishedAt: timestampMilli(r.PublishedAt),
		ExpiresAt:   timestampMilli(r.ExpiresAt),
		ImageURLs:   imageURLs[in.ID],
		Images:      images[in.ID],
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

func (s *sAnnouncement) SyncImages(ctx context.Context, in model.AnnouncementImageSyncInput) error {
	var current []announcementImageAttachment
	err := g.DB().Ctx(ctx).Model("attachments").
		Fields("id, announcement_id, file_id").
		Where("announcement_id = ?", in.AnnouncementID).
		Where("type = ?", consts.QuestionFileType).
		Order("id ASC").
		Scan(&current)
	if err != nil {
		return err
	}

	keep := make(map[int]struct{}, len(in.KeepFileIDs))
	for _, fileID := range in.KeepFileIDs {
		keep[fileID] = struct{}{}
	}
	deleteIDs := make([]int, 0, len(current))
	for _, attachment := range current {
		if _, ok := keep[attachment.FileId]; !ok {
			deleteIDs = append(deleteIDs, attachment.Id)
		}
	}
	if len(deleteIDs) > 0 {
		_, err = g.DB().Ctx(ctx).Model("attachments").WhereIn("id", deleteIDs).Delete()
		if err != nil {
			return err
		}
	}
	if len(in.AddFileIDs) == 0 {
		return nil
	}
	_, err = service.Attachment().AddAttachments(ctx, model.AddAttachmentInput{
		AnnouncementId: in.AnnouncementID,
		Type:           consts.QuestionFileType,
		FileId:         in.AddFileIDs,
	})
	return err
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
