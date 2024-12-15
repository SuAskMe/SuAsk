// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UpvotesDao is the data access object for table upvotes.
type UpvotesDao struct {
	table   string         // table is the underlying table name of the DAO.
	group   string         // group is the database configuration group name of current DAO.
	columns UpvotesColumns // columns contains all the column names of Table for convenient usage.
}

// UpvotesColumns defines and stores column names for table upvotes.
type UpvotesColumns struct {
	Id         string // 点赞ID
	UserId     string // 用户ID
	QuestionId string // 问题ID
	AnswerId   string // 回复ID
}

// upvotesColumns holds the columns for table upvotes.
var upvotesColumns = UpvotesColumns{
	Id:         "id",
	UserId:     "user_id",
	QuestionId: "question_id",
	AnswerId:   "answer_id",
}

// NewUpvotesDao creates and returns a new DAO object for table data access.
func NewUpvotesDao() *UpvotesDao {
	return &UpvotesDao{
		group:   "default",
		table:   "upvotes",
		columns: upvotesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *UpvotesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *UpvotesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *UpvotesDao) Columns() UpvotesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *UpvotesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *UpvotesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *UpvotesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
