// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserPoints is the golang structure of table user_points for DAO operations like Where/Data.
type UserPoints struct {
	g.Meta      `orm:"table:user_points, do:true"`
	Id          interface{} // ID
	UserId      interface{} // ID
	Points      interface{} //
	PointsTotal interface{} //
	CreatedAt   *gtime.Time //
	UpdatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
