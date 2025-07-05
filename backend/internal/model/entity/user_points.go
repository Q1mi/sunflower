// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPoints is the golang structure for table user_points.
type UserPoints struct {
	Id          uint64      `json:"id"          orm:"id"           description:"ID"` // ID
	UserId      uint64      `json:"userId"      orm:"user_id"      description:"ID"` // ID
	Points      int64       `json:"points"      orm:"points"       description:""`   //
	PointsTotal int64       `json:"pointsTotal" orm:"points_total" description:""`   //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""`   //
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:""`   //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""`   //
}
