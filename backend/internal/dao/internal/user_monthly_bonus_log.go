// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserMonthlyBonusLogDao is the data access object for the table user_monthly_bonus_log.
type UserMonthlyBonusLogDao struct {
	table    string                     // table is the underlying table name of the DAO.
	group    string                     // group is the database configuration group name of the current DAO.
	columns  UserMonthlyBonusLogColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler         // handlers for customized model modification.
}

// UserMonthlyBonusLogColumns defines and stores column names for the table user_monthly_bonus_log.
type UserMonthlyBonusLogColumns struct {
	Id          string // ID
	UserId      string // ID
	YearMonth   string // YYYYMM
	BonusType   string // 1:3 2:7 3:15 4:
	Description string //
	CreatedAt   string //
	UpdatedAt   string //
	DeletedAt   string //
}

// userMonthlyBonusLogColumns holds the columns for the table user_monthly_bonus_log.
var userMonthlyBonusLogColumns = UserMonthlyBonusLogColumns{
	Id:          "id",
	UserId:      "user_id",
	YearMonth:   "year_month",
	BonusType:   "bonus_type",
	Description: "description",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewUserMonthlyBonusLogDao creates and returns a new DAO object for table data access.
func NewUserMonthlyBonusLogDao(handlers ...gdb.ModelHandler) *UserMonthlyBonusLogDao {
	return &UserMonthlyBonusLogDao{
		group:    "default",
		table:    "user_monthly_bonus_log",
		columns:  userMonthlyBonusLogColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserMonthlyBonusLogDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserMonthlyBonusLogDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserMonthlyBonusLogDao) Columns() UserMonthlyBonusLogColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserMonthlyBonusLogDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserMonthlyBonusLogDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserMonthlyBonusLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
