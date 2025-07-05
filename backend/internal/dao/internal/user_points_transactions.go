// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserPointsTransactionsDao is the data access object for the table user_points_transactions.
type UserPointsTransactionsDao struct {
	table    string                        // table is the underlying table name of the DAO.
	group    string                        // group is the database configuration group name of the current DAO.
	columns  UserPointsTransactionsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler            // handlers for customized model modification.
}

// UserPointsTransactionsColumns defines and stores column names for the table user_points_transactions.
type UserPointsTransactionsColumns struct {
	Id              string // ID
	UserId          string // ID
	PointsChange    string // ()
	CurrentBalance  string //
	TransactionType string // (1: 2: 3: 4: 5:)
	Description     string //
	ExtJson         string //
	CreatedAt       string //
	UpdatedAt       string //
	DeletedAt       string //
}

// userPointsTransactionsColumns holds the columns for the table user_points_transactions.
var userPointsTransactionsColumns = UserPointsTransactionsColumns{
	Id:              "id",
	UserId:          "user_id",
	PointsChange:    "points_change",
	CurrentBalance:  "current_balance",
	TransactionType: "transaction_type",
	Description:     "description",
	ExtJson:         "ext_json",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
	DeletedAt:       "deleted_at",
}

// NewUserPointsTransactionsDao creates and returns a new DAO object for table data access.
func NewUserPointsTransactionsDao(handlers ...gdb.ModelHandler) *UserPointsTransactionsDao {
	return &UserPointsTransactionsDao{
		group:    "default",
		table:    "user_points_transactions",
		columns:  userPointsTransactionsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserPointsTransactionsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserPointsTransactionsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserPointsTransactionsDao) Columns() UserPointsTransactionsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserPointsTransactionsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserPointsTransactionsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserPointsTransactionsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
