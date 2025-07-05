// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Userinfo is the golang structure of table userinfo for DAO operations like Where/Data.
type Userinfo struct {
	g.Meta    `orm:"table:userinfo, do:true"`
	Id        interface{} // ID
	UserId    interface{} // ID
	Username  interface{} //
	Password  interface{} // (MD5)
	Email     interface{} //
	Avatar    interface{} //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
	DeletedAt *gtime.Time //
}
