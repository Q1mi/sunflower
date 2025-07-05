// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserMonthlyBonusLog is the golang structure for table user_monthly_bonus_log.
type UserMonthlyBonusLog struct {
	Id          uint64      `json:"id"          orm:"id"          description:"ID"`              // ID
	UserId      uint64      `json:"userId"      orm:"user_id"     description:"ID"`              // ID
	YearMonth   string      `json:"yearMonth"   orm:"year_month"  description:"YYYYMM"`          // YYYYMM
	BonusType   int         `json:"bonusType"   orm:"bonus_type"  description:"1:3 2:7 3:15 4:"` // 1:3 2:7 3:15 4:
	Description string      `json:"description" orm:"description" description:""`                //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:""`                //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:""`                //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"  description:""`                //
}
