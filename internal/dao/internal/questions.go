// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// QuestionsDao is the data access object for table questions.
type QuestionsDao struct {
	table   string           // table is the underlying table name of the DAO.
	group   string           // group is the database configuration group name of current DAO.
	columns QuestionsColumns // columns contains all the column names of Table for convenient usage.
}

// QuestionsColumns defines and stores column names for table questions.
type QuestionsColumns struct {
	Id        string // 问题ID
	SrcUserId string // 发起提问的用户ID
	DstUserId string // 被提问的用户ID，为空时问大家，不为空时问教师
	Title     string // 问题标题
	Contents  string // 问题内容
	IsPrivate string // 是否私密提问，仅在问教师时可为是
	CreatedAt string // 创建时间
	Views     string // 浏览量
	ReplyCnt  string // 回复数
}

// questionsColumns holds the columns for table questions.
var questionsColumns = QuestionsColumns{
	Id:        "id",
	SrcUserId: "src_user_id",
	DstUserId: "dst_user_id",
	Title:     "title",
	Contents:  "contents",
	IsPrivate: "is_private",
	CreatedAt: "created_at",
	Views:     "views",
	ReplyCnt:  "reply_cnt",
}

// NewQuestionsDao creates and returns a new DAO object for table data access.
func NewQuestionsDao() *QuestionsDao {
	return &QuestionsDao{
		group:   "default",
		table:   "questions",
		columns: questionsColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *QuestionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *QuestionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *QuestionsDao) Columns() QuestionsColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *QuestionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *QuestionsDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *QuestionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
