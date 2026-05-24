package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// ==================== 用户列表 ====================

type AdminUserItem struct {
	Id           int    `json:"id"            dc:"用户ID"`
	Name         string `json:"name"          dc:"用户名"`
	Nickname     string `json:"nickname"      dc:"昵称"`
	Email        string `json:"email"         dc:"邮箱"`
	Role         string `json:"role"          dc:"角色"`
	Introduction string `json:"introduction"  dc:"简介"`
	AvatarURL    string `json:"avatar"        dc:"头像URL"`
	CreatedAt    string `json:"created_at"    dc:"创建时间"`
}

type ListUsersReq struct {
	g.Meta  `path:"/admin/users" method:"GET" tags:"Admin" summary:"管理员-用户列表"`
	Page    int    `json:"page"    in:"query" v:"required|min:1" dc:"页码"`
	Role    string `json:"role"    in:"query" v:"in:admin,teacher,student,guest" dc:"角色筛选"`
	Keyword string `json:"keyword" in:"query" dc:"搜索关键词"`
}

type ListUsersRes struct {
	List       []AdminUserItem `json:"list"`
	Total      int             `json:"total"`
	RemainPage int             `json:"remain_page"`
}

// ==================== 创建用户 ====================

type CreateUserReq struct {
	g.Meta       `path:"/admin/users" method:"POST" tags:"Admin" summary:"管理员-创建用户"`
	Name         string `json:"name"         v:"required|length:1,50" dc:"用户名"`
	Password     string `json:"password"     v:"required|length:6,72" dc:"密码"`
	Email        string `json:"email"        v:"required|email" dc:"邮箱"`
	Role         string `json:"role"         v:"required|in:admin,teacher,student" dc:"角色"`
	Nickname     string `json:"nickname"     v:"required|length:1,50" dc:"昵称"`
	Introduction string `json:"introduction" dc:"教师简介(仅teacher)"`
	Perm         string `json:"perm"         dc:"提问箱权限(仅teacher): public/protected/private"`
}

type CreateUserRes struct {
	Id int `json:"id" dc:"新用户ID"`
}

// ==================== 编辑用户 ====================

type UpdateUserReq struct {
	g.Meta       `path:"/admin/users/{id}" method:"PUT" tags:"Admin" summary:"管理员-编辑用户"`
	Id           int    `json:"id"           in:"path" v:"required" dc:"用户ID"`
	Nickname     string `json:"nickname"     dc:"昵称"`
	Email        string `json:"email"        dc:"邮箱"`
	Role         string `json:"role"         dc:"角色"`
	Introduction string `json:"introduction" dc:"教师简介"`
	Perm         string `json:"perm"         dc:"提问箱权限"`
}

type UpdateUserRes struct {
	Id int `json:"id" dc:"用户ID"`
}

// ==================== 修改头像 ====================

type UpdateAvatarReq struct {
	g.Meta `path:"/admin/users/{id}/avatar" method:"PUT" mime:"multipart/form-data" tags:"Admin" summary:"管理员-修改用户头像"`
	Id     int `json:"id" in:"path" v:"required" dc:"用户ID"`
}

type UpdateAvatarRes struct {
	Id        int    `json:"id" dc:"用户ID"`
	AvatarURL string `json:"avatar" dc:"新头像URL"`
}

// ==================== 重置密码 ====================

type ResetPasswordReq struct {
	g.Meta   `path:"/admin/users/{id}/password" method:"PUT" tags:"Admin" summary:"管理员-重置密码"`
	Id       int    `json:"id"       in:"path" v:"required" dc:"用户ID"`
	Password string `json:"password" v:"required|length:6,64" dc:"新密码"`
}

type ResetPasswordRes struct {
	Id int `json:"id" dc:"用户ID"`
}

// ==================== 删除用户 ====================

type DeleteUserReq struct {
	g.Meta `path:"/admin/users/{id}" method:"DELETE" tags:"Admin" summary:"管理员-删除用户"`
	Id     int `json:"id" in:"path" v:"required" dc:"用户ID"`
}

type DeleteUserRes struct {
	Id int `json:"id" dc:"用户ID"`
}

// ==================== 内容管理：问题列表 ====================

type AdminQuestionItem struct {
	Id                 int    `json:"id"                   dc:"问题ID"`
	Title              string `json:"title"                dc:"问题标题"`
	Contents           string `json:"contents"             dc:"问题内容"`
	SrcUserId          int    `json:"src_user_id"          dc:"提问用户ID"`
	SrcUserName        string `json:"src_user_name"        dc:"提问用户名"`
	SrcUserNickname    string `json:"src_user_nickname"    dc:"提问用户昵称"`
	DstUserId          int    `json:"dst_user_id"          dc:"目标教师ID"`
	DstUserName        string `json:"dst_user_name"        dc:"目标教师用户名"`
	DstUserNickname    string `json:"dst_user_nickname"    dc:"目标教师昵称"`
	IsPrivate          bool   `json:"is_private"           dc:"是否私密"`
	CreatedAt          int64  `json:"created_at"           dc:"创建时间戳(毫秒)"`
	Views              int    `json:"views"                dc:"浏览量"`
	ReplyCnt           int    `json:"reply_cnt"            dc:"历史回复计数"`
	AnswerCount        int    `json:"answer_count"         dc:"当前可见回答数"`
	MatchedAnswerCount int    `json:"matched_answer_count" dc:"命中搜索词的回答数"`
	Status             string `json:"status"               dc:"answered/unanswered"`
	IsDeleted          bool   `json:"is_deleted"           dc:"是否已删除"`
	DeletedAt          int64  `json:"deleted_at"           dc:"删除时间戳(毫秒)"`
}

type ListQuestionsReq struct {
	g.Meta         `path:"/admin/questions" method:"GET" tags:"Admin" summary:"管理员-问题列表"`
	Page           int    `json:"page"             in:"query" v:"required|min:1" dc:"页码"`
	Keyword        string `json:"keyword"          in:"query" dc:"搜索关键词"`
	Status         string `json:"status"           in:"query" dc:"回答状态: all/answered/unanswered"`
	Visibility     string `json:"visibility"       in:"query" dc:"公开性: all/public/private"`
	TeacherId      int    `json:"teacher_id"       in:"query" dc:"目标教师ID"`
	IncludeDeleted bool   `json:"include_deleted" in:"query" dc:"是否包含已删除内容"`
}

type ListQuestionsRes struct {
	List       []AdminQuestionItem `json:"list"`
	Total      int                 `json:"total"`
	RemainPage int                 `json:"remain_page"`
}

// ==================== 内容管理：问题详情与回答 ====================

type AdminQuestionAnswerItem struct {
	Id           int    `json:"id"            dc:"回答ID"`
	QuestionId   int    `json:"question_id"   dc:"问题ID"`
	UserId       int    `json:"user_id"       dc:"回答用户ID"`
	UserName     string `json:"user_name"     dc:"回答用户名"`
	UserNickname string `json:"user_nickname" dc:"回答用户昵称"`
	UserRole     string `json:"user_role"     dc:"回答用户角色"`
	Contents     string `json:"contents"      dc:"回答内容"`
	CreatedAt    int64  `json:"created_at"    dc:"创建时间戳(毫秒)"`
	Upvotes      int    `json:"upvotes"       dc:"点赞数"`
	InReplyTo    int    `json:"in_reply_to"   dc:"回复的回答ID"`
	IsDeleted    bool   `json:"is_deleted"    dc:"是否已删除"`
	DeletedAt    int64  `json:"deleted_at"    dc:"删除时间戳(毫秒)"`
}

type GetQuestionDetailReq struct {
	g.Meta         `path:"/admin/questions/{id}" method:"GET" tags:"Admin" summary:"管理员-问题详情"`
	Id             int  `json:"id"              in:"path"  v:"required|min:1" dc:"问题ID"`
	IncludeDeleted bool `json:"include_deleted" in:"query" dc:"是否包含已删除回答"`
}

type GetQuestionDetailRes struct {
	Question AdminQuestionItem         `json:"question"`
	Answers  []AdminQuestionAnswerItem `json:"answers"`
}

// ==================== 内容管理：删除问题/回答 ====================

type DeleteQuestionReq struct {
	g.Meta `path:"/admin/questions/{id}" method:"DELETE" tags:"Admin" summary:"管理员-删除问题"`
	Id     int `json:"id" in:"path" v:"required|min:1" dc:"问题ID"`
}

type DeleteQuestionRes struct {
	Id int `json:"id" dc:"问题ID"`
}

type DeleteQuestionAnswerReq struct {
	g.Meta     `path:"/admin/questions/{question_id}/answers/{answer_id}" method:"DELETE" tags:"Admin" summary:"管理员-删除问题下的回答"`
	QuestionId int `json:"question_id" in:"path" v:"required|min:1" dc:"问题ID"`
	AnswerId   int `json:"answer_id"   in:"path" v:"required|min:1" dc:"回答ID"`
}

type DeleteQuestionAnswerRes struct {
	Id         int `json:"id"          dc:"回答ID"`
	QuestionId int `json:"question_id" dc:"问题ID"`
}
