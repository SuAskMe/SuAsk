package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// --- 活跃公告（Banner） ---

type ActiveReq struct {
	g.Meta `path:"/announcements/active" method:"GET" tags:"Announcement" summary:"当前活跃公告"`
}

type ActiveRes struct {
	Announcement *ActiveAnnouncementItem `json:"announcement"`
}

type ActiveAnnouncementItem struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"contents"`
	AuthorName  string   `json:"author_name"`
	IsPinned    bool     `json:"is_pinned"`
	PublishedAt int64    `json:"published_at"`
	ExpiresAt   int64    `json:"expires_at"`
	ImageURLs   []string `json:"image_urls"`
}

// --- 列表 ---

type ListReq struct {
	g.Meta `path:"/announcements" method:"GET" tags:"Announcement" summary:"公告列表"`
	Page   int `json:"page" v:"required|min:1"`
}

type ListRes struct {
	Announcements []AnnouncementItem `json:"announcements"`
	RemainPage    int                `json:"remain_page"`
	Total         int                `json:"total"`
}

type AdminListReq struct {
	g.Meta `path:"/announcements/admin" method:"GET" tags:"Announcement" summary:"公告管理列表（仅管理员）"`
	Page   int `json:"page" v:"required|min:1"`
}

type AdminListRes = ListRes

type AnnouncementItem struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Content     string   `json:"contents"`
	AuthorName  string   `json:"author_name"`
	IsPinned    bool     `json:"is_pinned"`
	PublishedAt int64    `json:"published_at"`
	ExpiresAt   int64    `json:"expires_at"`
	CommentCnt  int      `json:"comment_cnt"`
	ImageURLs   []string `json:"image_urls"`
}

type AnnouncementImage struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// --- 详情 ---

type DetailReq struct {
	g.Meta `path:"/announcements/detail" method:"GET" tags:"Announcement" summary:"公告详情"`
	ID     int `json:"id" v:"required|min:1"`
}

type DetailRes struct {
	ID          int                 `json:"id"`
	Title       string              `json:"title"`
	Content     string              `json:"contents"`
	AuthorName  string              `json:"author_name"`
	IsPinned    bool                `json:"is_pinned"`
	PublishedAt int64               `json:"published_at"`
	ExpiresAt   int64               `json:"expires_at"`
	ImageURLs   []string            `json:"image_urls"`
	Images      []AnnouncementImage `json:"images"`
	Comments    []CommentItem       `json:"comments"`
}

type CommentItem struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Nickname  string `json:"nickname"`
	Avatar    string `json:"avatar"`
	Contents  string `json:"contents"`
	CreatedAt int64  `json:"created_at"`
	InReplyTo int    `json:"in_reply_to"`
}

// --- 发布 ---

type CreateReq struct {
	g.Meta    `path:"/announcements" method:"POST" tags:"Announcement" summary:"发布公告（仅管理员）"`
	Title     string              `json:"title" v:"required"`
	Content   string              `json:"content" v:"required"`
	IsPinned  bool                `json:"is_pinned"`
	ExpiresAt string              `json:"expires_at"` // ISO8601 或空
	Files     []*ghttp.UploadFile `json:"files"`
}

type CreateRes struct {
	ID int `json:"id"`
}

// --- 编辑 ---

type UpdateReq struct {
	g.Meta       `path:"/announcements" method:"PUT" tags:"Announcement" summary:"编辑公告（仅管理员）"`
	ID           int                 `json:"id" v:"required|min:1"`
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	IsPinned     *bool               `json:"is_pinned"`
	ExpiresAt    *string             `json:"expires_at"`
	SyncImages   bool                `json:"sync_images"`
	KeepImageIDs string              `json:"keep_image_ids"`
	Files        []*ghttp.UploadFile `json:"files"`
}

type UpdateRes struct {
	ID int `json:"id"`
}

// --- 删除（软删） ---

type DeleteReq struct {
	g.Meta `path:"/announcements" method:"DELETE" tags:"Announcement" summary:"删除公告（仅管理员）"`
	ID     int `json:"id" v:"required|min:1"`
}

type DeleteRes struct{}

// --- 评论 ---

type AddCommentReq struct {
	g.Meta         `path:"/announcements/comment" method:"POST" tags:"Announcement" summary:"评论公告"`
	AnnouncementID int    `json:"announcement_id" v:"required|min:1"`
	Content        string `json:"content" v:"required"`
	InReplyTo      *int   `json:"in_reply_to"`
}

type AddCommentRes struct {
	ID int `json:"id"`
}
