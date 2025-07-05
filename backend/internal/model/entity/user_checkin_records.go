// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserCheckinRecords is the golang structure for table user_checkin_records.
type UserCheckinRecords struct {
	Id                uint64      `json:"id"                orm:"id"                  description:"ID"`   // ID
	UserId            uint64      `json:"userId"            orm:"user_id"             description:"ID"`   // ID
	CheckinDate       *gtime.Time `json:"checkinDate"       orm:"checkin_date"        description:""`     //
	CheckinType       int         `json:"checkinType"       orm:"checkin_type"        description:"1=2="` // 1=2=
	PointsAwardedBase int         `json:"pointsAwardedBase" orm:"points_awarded_base" description:""`     //
	CreatedAt         *gtime.Time `json:"createdAt"         orm:"created_at"          description:""`     //
	UpdatedAt         *gtime.Time `json:"updatedAt"         orm:"updated_at"          description:""`     //
	DeletedAt         *gtime.Time `json:"deletedAt"         orm:"deleted_at"          description:""`     //
}
