// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPointsTransactions is the golang structure for table user_points_transactions.
type UserPointsTransactions struct {
	Id              uint64      `json:"id"              orm:"id"               description:"ID"`               // ID
	UserId          uint64      `json:"userId"          orm:"user_id"          description:"ID"`               // ID
	PointsChange    int64       `json:"pointsChange"    orm:"points_change"    description:"()"`               // ()
	CurrentBalance  int64       `json:"currentBalance"  orm:"current_balance"  description:""`                 //
	TransactionType int         `json:"transactionType" orm:"transaction_type" description:"(1: 2: 3: 4: 5:)"` // (1: 2: 3: 4: 5:)
	Description     string      `json:"description"     orm:"description"      description:""`                 //
	ExtJson         string      `json:"extJson"         orm:"ext_json"         description:""`                 //
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:""`                 //
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:""`                 //
	DeletedAt       *gtime.Time `json:"deletedAt"       orm:"deleted_at"       description:""`                 //
}
