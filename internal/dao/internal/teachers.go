// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// TeachersDao is the data access object for the table teachers.
type TeachersDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  TeachersColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// TeachersColumns defines and stores column names for the table teachers.
type TeachersColumns struct {
	Id           string //
	Responses    string // 回复数
	Name         string // 老师名字
	AvatarUrl    string // 老师头像链接
	Introduction string // 老师简介
	Email        string // 老师邮箱
	Perm         string // 提问箱权限
}

// teachersColumns holds the columns for the table teachers.
var teachersColumns = TeachersColumns{
	Id:           "id",
	Responses:    "responses",
	Name:         "name",
	AvatarUrl:    "avatar_url",
	Introduction: "introduction",
	Email:        "email",
	Perm:         "perm",
}

// NewTeachersDao creates and returns a new DAO object for table data access.
func NewTeachersDao(handlers ...gdb.ModelHandler) *TeachersDao {
	return &TeachersDao{
		group:    "default",
		table:    "teachers",
		columns:  teachersColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *TeachersDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *TeachersDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *TeachersDao) Columns() TeachersColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *TeachersDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *TeachersDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *TeachersDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
