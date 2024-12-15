// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Attachments is the golang structure of table attachments for DAO operations like Where/Data.
type Attachments struct {
	g.Meta     `orm:"table:attachments, do:true"`
	Id         interface{} // 附件ID
	QuestionId interface{} // 问题ID
	AnswerId   interface{} // 回答ID
	Type       interface{} // 附件类型（目前仅支持图片）
	FileId     interface{} // 文件ID
}
