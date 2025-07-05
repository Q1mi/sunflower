// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Userinfo is the golang structure for table userinfo.
type Userinfo struct {
	Id        uint64      `json:"id"        orm:"id"         description:"ID"`    // ID
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"ID"`    // ID
	Username  string      `json:"username"  orm:"username"   description:""`      //
	Password  string      `json:"password"  orm:"password"   description:"(MD5)"` // (MD5)
	Email     string      `json:"email"     orm:"email"      description:""`      //
	Avatar    string      `json:"avatar"    orm:"avatar"     description:""`      //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`      //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`      //
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:""`      //
}
