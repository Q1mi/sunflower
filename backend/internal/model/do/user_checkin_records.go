// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCheckinRecords is the golang structure of table user_checkin_records for DAO operations like Where/Data.
type UserCheckinRecords struct {
	g.Meta            `orm:"table:user_checkin_records, do:true"`
	Id                interface{} // ID
	UserId            interface{} // ID
	CheckinDate       *gtime.Time //
	CheckinType       interface{} // 1=2=
	PointsAwardedBase interface{} //
	CreatedAt         *gtime.Time //
	UpdatedAt         *gtime.Time //
	DeletedAt         *gtime.Time //
}
