// =================================================================================
// This file is auto-generated by the GoFrame CLI tool. You may modify it as needed.
// =================================================================================

package dao

import (
	"backend/internal/dao/internal"
)

// userMonthlyBonusLogDao is the data access object for the table user_monthly_bonus_log.
// You can define custom methods on it to extend its functionality as needed.
type userMonthlyBonusLogDao struct {
	*internal.UserMonthlyBonusLogDao
}

var (
	// UserMonthlyBonusLog is a globally accessible object for table user_monthly_bonus_log operations.
	UserMonthlyBonusLog = userMonthlyBonusLogDao{internal.NewUserMonthlyBonusLogDao()}
)

// Add your custom methods and functionality below.
