// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnswersDao is the data access object for table answers.
type AnswersDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns AnswersColumns // columns contains all the column names of Table for convenient usage.
}

// AnswersColumns defines and stores column names for table answers.
type AnswersColumns struct {
	Id         string // 回答ID
	UserId     string // 用户ID
	QuestionId string // 问题ID
	InReplyTo  string // 回复的回答ID，可为空
	Contents   string // 回答内容
	CreatedAt  string // 创建时间
	Upvotes    string // 点赞量
}

// answersColumns holds the columns for table answers.
var answersColumns = AnswersColumns{
	Id:         "id",
	UserId:     "user_id",
	QuestionId: "question_id",
	InReplyTo:  "in_reply_to",
	Contents:   "contents",
	CreatedAt:  "created_at",
	Upvotes:    "upvotes",
}

// NewAnswersDao creates and returns a new DAO object for table data access.
func NewAnswersDao() *AnswersDao {
	return &AnswersDao{
		group:   "default",
		table:   "answers",
		columns: answersColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *AnswersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *AnswersDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *AnswersDao) Columns() AnswersColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *AnswersDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *AnswersDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *AnswersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
