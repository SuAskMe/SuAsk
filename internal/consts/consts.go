package consts

const (
	DefaultAvatarURL      = "default-avatar"
	NumOfQuestionsPerPage = 30
	SortByTimeDsc         = 0
	SortByTimeAsc         = 1
	SortByViewsDsc        = 2
	SortByViewsAsc        = 3
)

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
	TokenType       = "Bearer"
	CacheMode       = 1 // gcache
	ServerName      = "SuAsk"
	ErrLoginFaulMsg = "登录失败，账号或密码错误"
)

// for login
const (
	CtxId   = "CtxId"
	CtxName = "CtxName"
	CtxRole = "CtxRole"
)

// for register
const (
	DefaultThemeId = 1
)

// for file
const (
	FileUploadMaxMinutes = 10
	FileServerPrefix     = "http://localhost:8080"
)

// for question

const DefaultUserId = 1
const QuestionFileType = "picture"
