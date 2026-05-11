package announcement

import (
	"context"
	"fmt"
	v1 "suask/api/announcement/v1"
	"suask/internal/consts"
	"suask/internal/model"
	"suask/internal/service"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type cAnnouncement struct{}

var Announcement = cAnnouncement{}

// --- 列表（所有人可访问） ---

func (c *cAnnouncement) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := service.Announcement().List(ctx, model.AnnouncementListInput{Page: req.Page})
	if err != nil {
		return nil, err
	}
	items := make([]v1.AnnouncementItem, len(out.Items))
	for i, it := range out.Items {
		items[i] = v1.AnnouncementItem{
			ID:          it.ID,
			Title:       it.Title,
			Content:     it.Content,
			AuthorName:  it.AuthorName,
			IsPinned:    it.IsPinned,
			PublishedAt: it.PublishedAt,
			CommentCnt:  it.CommentCnt,
		}
	}
	return &v1.ListRes{Announcements: items, RemainPage: out.RemainPage}, nil
}

// --- 详情（所有人可访问） ---

func (c *cAnnouncement) Detail(ctx context.Context, req *v1.DetailReq) (res *v1.DetailRes, err error) {
	out, err := service.Announcement().Detail(ctx, model.AnnouncementDetailInput{ID: req.ID})
	if err != nil {
		return nil, err
	}
	// 图片 URL
	var imageURLs []string
	if len(out.ImageIDs) > 0 {
		fileList, err := service.File().GetList(ctx, model.FileListGetInput{IdList: out.ImageIDs})
		if err == nil {
			imageURLs = fileList.URL
		}
	}
	// 评论
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
		ImageURLs:   imageURLs,
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
	var expiresAt *gtime.Time
	if req.ExpiresAt != "" {
		expiresAt = gtime.NewFromStr(req.ExpiresAt)
	}
	out, err := service.Announcement().Update(ctx, model.AnnouncementUpdateInput{
		ID:        req.ID,
		Title:     req.Title,
		Content:   req.Content,
		IsPinned:  req.IsPinned,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return nil, err
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
