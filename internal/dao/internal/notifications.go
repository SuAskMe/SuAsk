// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NotificationsDao is the data access object for the table notifications.
type NotificationsDao struct {
	table    string               // table is the underlying table name of the DAO.
	group    string               // group is the database configuration group name of the current DAO.
	columns  NotificationsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler   // handlers for customized model modification.
}

// NotificationsColumns defines and stores column names for the table notifications.
type NotificationsColumns struct {
	Id         string // 提醒ID
	UserId     string // 用户ID
	QuestionId string // 问题ID
	ReplyToId  string // 回复问题的ID
	AnswerId   string // 问题ID
	Type       string // 提醒类型（新提问、新回复、新回答）
	IsRead     string // 是否已读
	CreatedAt  string // 创建时间
	DeletedAt  string // 删除时间
}

// notificationsColumns holds the columns for the table notifications.
var notificationsColumns = NotificationsColumns{
	Id:         "id",
	UserId:     "user_id",
	QuestionId: "question_id",
	ReplyToId:  "reply_to_id",
	AnswerId:   "answer_id",
	Type:       "type",
	IsRead:     "is_read",
	CreatedAt:  "created_at",
	DeletedAt:  "deleted_at",
}

// NewNotificationsDao creates and returns a new DAO object for table data access.
func NewNotificationsDao(handlers ...gdb.ModelHandler) *NotificationsDao {
	return &NotificationsDao{
		group:    "default",
		table:    "notifications",
		columns:  notificationsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NotificationsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NotificationsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NotificationsDao) Columns() NotificationsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NotificationsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NotificationsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *NotificationsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
