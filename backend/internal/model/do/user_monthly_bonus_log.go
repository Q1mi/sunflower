// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserMonthlyBonusLog is the golang structure of table user_monthly_bonus_log for DAO operations like Where/Data.
type UserMonthlyBonusLog struct {
	g.Meta      `orm:"table:user_monthly_bonus_log, do:true"`
	Id          interface{} // ID
	UserId      interface{} // ID
	YearMonth   interface{} // YYYYMM
	BonusType   interface{} // 1:3 2:7 3:15 4:
	Description interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
