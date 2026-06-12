// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Questions is the golang structure of table questions for DAO operations like Where/Data.
type Questions struct {
	g.Meta    `orm:"table:questions, do:true"`
	Id        any         //
	SrcUserId any         //
	DstUserId any         //
	Title     any         //
	Contents  any         //
	CreatedAt *gtime.Time //
	Views     any         //
	ReplyCnt  any         //
	DeletedAt *gtime.Time //
}
