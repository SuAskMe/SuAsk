package consts

const (
	MaxQuestionsPerPage   = 10
	MaxKeywordsPerReq     = 8
	MaxAvatarsPerQuestion = 3
	SortByTimeDsc         = 0
	SortByTimeAsc         = 1
	SortByViewsDsc        = 2
	SortByViewsAsc        = 3
)

// default settings

const (
	DefaultUserId    = 1
	DefaultUserName  = "匿名用户"
	DefaultThemeId   = 1
	DefaultAvatarURL = "default-avatar"
	ErrInternal      = "出错了，请联系管理员"
)

// for question status

const (
	Answered   = "已回答"
	Unanswered = "未回答"
	OnTop      = "置顶"
)

// for User Role

const (
	TEACHER = "teacher"
	STUDENT = "student"
	ADMIN   = "admin"
)

// for gToken
const (
	TokenType  = "Bearer"
	ServerName = "SuAsk"
)

// for login
const (
	CtxId   = "CtxId"
	CtxName = "CtxName"
	CtxRole = "CtxRole"
)

// for question

const QuestionFileType = "picture"

// for notification
const (
	NewQuestion = "new_question"
	NewAnswer   = "new_answer"
	NewReply    = "new_reply"
)

const (
	ForgetPassword = "forget_password"
	ResetPassword  = "reset_password"
)

const (
	PermPublic    = "public"
	PermProtected = "protected"
	PermPrivate   = "private"
)
