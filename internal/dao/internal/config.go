// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// ConfigDao is the data access object for the table config.
type ConfigDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  ConfigColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// ConfigColumns defines and stores column names for the table config.
type ConfigColumns struct {
	Id                string // 配置ID，限制为0
	DefaultAvatarPath string // 默认头像文件路径
	DefaultThemeId    string // 默认主题ID
}

// configColumns holds the columns for the table config.
var configColumns = ConfigColumns{
	Id:                "id",
	DefaultAvatarPath: "default_avatar_path",
	DefaultThemeId:    "default_theme_id",
}

// NewConfigDao creates and returns a new DAO object for table data access.
func NewConfigDao(handlers ...gdb.ModelHandler) *ConfigDao {
	return &ConfigDao{
		group:    "default",
		table:    "config",
		columns:  configColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *ConfigDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *ConfigDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *ConfigDao) Columns() ConfigColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *ConfigDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *ConfigDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *ConfigDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
