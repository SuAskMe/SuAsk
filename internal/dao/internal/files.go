// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FilesDao is the data access object for table files.
type FilesDao struct {
	table   string       // table is the underlying table name of the DAO.
	group   string       // group is the database configuration group name of current DAO.
	columns FilesColumns // columns contains all the column names of Table for convenient usage.
}

// FilesColumns defines and stores column names for table files.
type FilesColumns struct {
	Id         string // 文件ID
	Name       string // 文件名，不得包含非法字符例如斜杠
	Hash       string // 文件哈希，算法暂定为BLAKE2b
	UploaderId string // 上传者用户ID
	CreatedAt  string //
}

// filesColumns holds the columns for table files.
var filesColumns = FilesColumns{
	Id:         "id",
	Name:       "name",
	Hash:       "hash",
	UploaderId: "uploader_id",
	CreatedAt:  "created_at",
}

// NewFilesDao creates and returns a new DAO object for table data access.
func NewFilesDao() *FilesDao {
	return &FilesDao{
		group:   "default",
		table:   "files",
		columns: filesColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *FilesDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of current dao.
func (dao *FilesDao) Table() string {
	return dao.table
}

// Columns returns all column names of current dao.
func (dao *FilesDao) Columns() FilesColumns {
	return dao.columns
}

// Group returns the configuration group name of database of current dao.
func (dao *FilesDao) Group() string {
	return dao.group
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *FilesDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *FilesDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
