package announcement

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	v1 "suask/api/announcement/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type cAnnouncement struct{}

var Announcement = cAnnouncement{}

func toAnnouncementImage(it model.AnnouncementImage) v1.AnnouncementImage {
	return v1.AnnouncementImage{
		ID:  it.ID,
		URL: it.URL,
	}
}

func toAnnouncementImages(items []model.AnnouncementImage) []v1.AnnouncementImage {
	res := make([]v1.AnnouncementImage, len(items))
	for i, it := range items {
		res[i] = toAnnouncementImage(it)
	}
	return res
}

func toAnnouncementItem(it model.AnnouncementListItem) v1.AnnouncementItem {
	return v1.AnnouncementItem{
		ID:          it.ID,
		Title:       it.Title,
		Content:     it.Content,
		AuthorName:  it.AuthorName,
		IsPinned:    it.IsPinned,
		PublishedAt: it.PublishedAt,
		ExpiresAt:   it.ExpiresAt,
		CommentCnt:  it.CommentCnt,
		ImageURLs:   it.ImageURLs,
	}
}

func toActiveAnnouncementItem(it model.AnnouncementActiveItem) v1.ActiveAnnouncementItem {
	return v1.ActiveAnnouncementItem{
		ID:          it.ID,
		Title:       it.Title,
		Content:     it.Content,
		AuthorName:  it.AuthorName,
		IsPinned:    it.IsPinned,
		PublishedAt: it.PublishedAt,
		ExpiresAt:   it.ExpiresAt,
		ImageURLs:   it.ImageURLs,
	}
}

func toAnnouncementItems(items []model.AnnouncementListItem) []v1.AnnouncementItem {
	res := make([]v1.AnnouncementItem, len(items))
	for i, it := range items {
		res[i] = toAnnouncementItem(it)
	}
	return res
}

func toListRes(out *model.AnnouncementListOutput) *v1.ListRes {
	return &v1.ListRes{
		Announcements: toAnnouncementItems(out.Items),
		RemainPage:    out.RemainPage,
		Total:         out.Total,
	}
}

func parseKeepImageIDs(raw string) []int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	var ids []int
	if strings.HasPrefix(raw, "[") {
		if err := json.Unmarshal([]byte(raw), &ids); err == nil {
			return ids
		}
	}
	parts := strings.Split(raw, ",")
	ids = make([]int, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err == nil && id > 0 {
			ids = append(ids, id)
		}
	}
	return ids
}

// --- 列表（所有人可访问） ---

func (c *cAnnouncement) GetActive(ctx context.Context, req *v1.ActiveReq) (res *v1.ActiveRes, err error) {
	out, err := service.Announcement().GetActive(ctx)
	if err != nil {
		return nil, err
	}
	res = &v1.ActiveRes{}
	if out.Item != nil {
		item := toActiveAnnouncementItem(*out.Item)
		res.Announcement = &item
	}
	return res, nil
}

// --- 列表（所有人可访问） ---

func (c *cAnnouncement) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := service.Announcement().List(ctx, model.AnnouncementListInput{Page: req.Page})
	if err != nil {
		return nil, err
	}
	return toListRes(out), nil
}

func (c *cAnnouncement) AdminList(ctx context.Context, req *v1.AdminListReq) (res *v1.AdminListRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if err := requireAdmin(ctx, uid); err != nil {
		return nil, err
	}
	out, err := service.Announcement().List(ctx, model.AnnouncementListInput{Page: req.Page, IncludeExpired: true})
	if err != nil {
		return nil, err
	}
	return toListRes(out), nil
}

// --- 详情（所有人可访问） ---

func (c *cAnnouncement) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	out, err := service.Announcement().Detail(ctx, model.AnnouncementDetailInput{ID: req.ID})
	if err != nil {
		return nil, err
	}
	comments, _ := service.Announcement().GetComments(ctx, req.ID)
	commentItems := make([]v1.CommentItem, len(comments))
	for i, c := range comments {
		commentItems[i] = v1.CommentItem{
			ID:        c.ID,
			UserID:    c.UserID,
			Nickname:  c.Nickname,
			Avatar:    c.AvatarURL,
			Contents:  c.Contents,
			CreatedAt: c.CreatedAt,
			InReplyTo: c.InReplyTo,
		}
	}
	return &v1.DetailRes{
		ID:          out.ID,
		Title:       out.Title,
		Content:     out.Content,
		AuthorName:  out.AuthorName,
		IsPinned:    out.IsPinned,
		PublishedAt: out.PublishedAt,
		ExpiresAt:   out.ExpiresAt,
		ImageURLs:   out.ImageURLs,
		Images:      toAnnouncementImages(out.Images),
		Comments:    commentItems,
	}, nil
}

// --- 发布（仅 admin） ---

func (c *cAnnouncement) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if err := requireAdmin(ctx, uid); err != nil {
		return nil, err
	}
	var expiresAt *gtime.Time
	if req.ExpiresAt != "" {
		expiresAt = gtime.NewFromStr(req.ExpiresAt)
	}
	out, err := service.Announcement().Create(ctx, model.AnnouncementCreateInput{
		AuthorID:  uid,
		Title:     req.Title,
		Content:   req.Content,
		IsPinned:  req.IsPinned,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
	}
	// 附件
	if req.Files != nil && len(req.Files) > 0 {
		fileList := model.FileListAddInput{FileList: req.Files, UploaderId: uid}
		fileIdList, err := service.File().UploadFileList(ctx, fileList)
		if err != nil {
			return nil, err
		}
		attachment := model.AddAttachmentInput{
			AnnouncementId: out.ID,
			Type:           consts.QuestionFileType,
			FileId:         fileIdList.IdList,
		}
		_, err = service.Attachment().AddAttachments(ctx, attachment)
		if err != nil {
			return nil, err
		}
	}
	return &v1.CreateRes{ID: out.ID}, nil
}

// --- 编辑（仅 admin） ---

func (c *cAnnouncement) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if err := requireAdmin(ctx, uid); err != nil {
		return nil, err
	}
	var (
		expiresAt      *gtime.Time
		clearExpiresAt bool
	)
	if req.ExpiresAt != nil {
		if *req.ExpiresAt == "" {
			clearExpiresAt = true
		} else {
			expiresAt = gtime.NewFromStr(*req.ExpiresAt)
		}
	}
	out, err := service.Announcement().Update(ctx, model.AnnouncementUpdateInput{
		ID:             req.ID,
		Title:          req.Title,
		Content:        req.Content,
		IsPinned:       req.IsPinned,
		ExpiresAt:      expiresAt,
		ClearExpiresAt: clearExpiresAt,
	})
	if err != nil {
		return nil, err
	}
	if req.SyncImages {
		var addFileIDs []int
		if len(req.Files) > 0 {
			fileList := model.FileListAddInput{FileList: req.Files, UploaderId: uid}
			fileIdList, err := service.File().UploadFileList(ctx, fileList)
			if err != nil {
				return nil, err
			}
			addFileIDs = fileIdList.IdList
		}
		if err := service.Announcement().SyncImages(ctx, model.AnnouncementImageSyncInput{
			AnnouncementID: req.ID,
			KeepFileIDs:    parseKeepImageIDs(req.KeepImageIDs),
			AddFileIDs:     addFileIDs,
		}); err != nil {
			return nil, err
		}
	}
	return &v1.UpdateRes{ID: out.ID}, nil
}

// --- 删除（仅 admin） ---

func (c *cAnnouncement) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if err := requireAdmin(ctx, uid); err != nil {
		return nil, err
	}
	if err := service.Announcement().Delete(ctx, model.AnnouncementDeleteInput{ID: req.ID}); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}

// --- 评论（登录用户） ---

func (c *cAnnouncement) AddComment(ctx context.Context, req *v1.AddCommentReq) (res *v1.AddCommentRes, err error) {
	uid := gconv.Int(ctx.Value(consts.CtxId))
	if uid == consts.DefaultUserId {
		return nil, fmt.Errorf("请登录后再评论")
	}
	var inReplyTo *int
	if req.InReplyTo != nil {
		inReplyTo = req.InReplyTo
	}
	out, err := service.Announcement().AddComment(ctx, model.AnnouncementCommentInput{
		AnnouncementID: req.AnnouncementID,
		UserID:         uid,
		Content:        req.Content,
		InReplyTo:      inReplyTo,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AddCommentRes{ID: out.ID}, nil
}

// --- 权限 helper ---

func requireAdmin(ctx context.Context, uid int) error {
	if uid == consts.DefaultUserId {
		return fmt.Errorf("请登录")
	}
	user, err := service.User().GetUser(ctx, model.UserInfoInput{Id: uid})
	if err != nil {
		return err
	}
	if user.Role != consts.ADMIN {
		return fmt.Errorf("仅管理员可操作")
	}
	return nil
}
