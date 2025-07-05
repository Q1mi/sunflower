// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// UserCheckinRecordsDao is the data access object for the table user_checkin_records.
type UserCheckinRecordsDao struct {
	table    string                    // table is the underlying table name of the DAO.
	group    string                    // group is the database configuration group name of the current DAO.
	columns  UserCheckinRecordsColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler        // handlers for customized model modification.
}

// UserCheckinRecordsColumns defines and stores column names for the table user_checkin_records.
type UserCheckinRecordsColumns struct {
	Id                string // ID
	UserId            string // ID
	CheckinDate       string //
	CheckinType       string // 1=2=
	PointsAwardedBase string //
	CreatedAt         string //
	UpdatedAt         string //
	DeletedAt         string //
}

// userCheckinRecordsColumns holds the columns for the table user_checkin_records.
var userCheckinRecordsColumns = UserCheckinRecordsColumns{
	Id:                "id",
	UserId:            "user_id",
	CheckinDate:       "checkin_date",
	CheckinType:       "checkin_type",
	PointsAwardedBase: "points_awarded_base",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	DeletedAt:         "deleted_at",
}

// NewUserCheckinRecordsDao creates and returns a new DAO object for table data access.
func NewUserCheckinRecordsDao(handlers ...gdb.ModelHandler) *UserCheckinRecordsDao {
	return &UserCheckinRecordsDao{
		group:    "default",
		table:    "user_checkin_records",
		columns:  userCheckinRecordsColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *UserCheckinRecordsDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *UserCheckinRecordsDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *UserCheckinRecordsDao) Columns() UserCheckinRecordsColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *UserCheckinRecordsDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *UserCheckinRecordsDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *UserCheckinRecordsDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
