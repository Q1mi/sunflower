// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPointsTransactions is the golang structure of table user_points_transactions for DAO operations like Where/Data.
type UserPointsTransactions struct {
	g.Meta          `orm:"table:user_points_transactions, do:true"`
	Id              interface{} // ID
	UserId          interface{} // ID
	PointsChange    interface{} // ()
	CurrentBalance  interface{} //
	TransactionType interface{} // (1: 2: 3: 4: 5:)
	Description     interface{} //
	ExtJson         interface{} //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time //
}
