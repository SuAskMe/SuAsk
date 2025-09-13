// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SettingsDao is the data access object for the table settings.
type SettingsDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SettingsColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SettingsColumns defines and stores column names for the table settings.
type SettingsColumns struct {
	Id              string // 设置id
	ThemeId         string // 主题ID，为空时为配置的默认主题
	QuestionBoxPerm string // 提问箱权限
	NotifyEmail     string // 通知邮箱
	NotifyMergeCnt  string // 合并提醒数量
	NotifyMaxDelay  string // 最大延迟时间(min)
}

// settingsColumns holds the columns for the table settings.
var settingsColumns = SettingsColumns{
	Id:              "id",
	ThemeId:         "theme_id",
	QuestionBoxPerm: "question_box_perm",
	NotifyEmail:     "notify_email",
	NotifyMergeCnt:  "notify_merge_cnt",
	NotifyMaxDelay:  "notify_max_delay",
}

// NewSettingsDao creates and returns a new DAO object for table data access.
func NewSettingsDao(handlers ...gdb.ModelHandler) *SettingsDao {
	return &SettingsDao{
		group:    "default",
		table:    "settings",
		columns:  settingsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SettingsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SettingsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SettingsDao) Columns() SettingsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SettingsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SettingsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *SettingsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
