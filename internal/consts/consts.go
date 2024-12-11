package consts


const (
	DefaultAvatarURL      = "/src/assets/default_avatar.png"
	NumOfQuestionsPerPage = 30
	SortByTimeDsc         = 0
	SortByTimeAsc         = 1
	SortByViewsDsc        = 2
	SortByViewsAsc        = 3
)
// for User Role

const TEACHER = "teacher"
const STUDENT = "student"
const ADMIN = "admin"

// for gToken

const TokenType = "Bearer"
const CacheMode = 1 // gcache
const ServerName = "SuAsk"
const ErrLoginFaulMsg = "登录失败，账号或密码错误"

const (
	CtxId   = "CtxId"
	CtxName = "CtxName"
	CtxRole = "CtxRole"
)

// const DefaultAvatarURL = "https://p1-arco.byteimg.com/tos-cn-i-uwbnlip3yd/3ee5f13fb09879ecb5185e440cef6eb9.png~tplv-uwbnlip3yd-webp.webp"
const DefaultThemeId = 0

const FileUploadMaxMinutes = 10
