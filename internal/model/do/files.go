// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Files is the golang structure of table files for DAO operations like Where/Data.
type Files struct {
	g.Meta     `orm:"table:files, do:true"`
	Id         interface{} // 文件ID
	Name       interface{} // 文件名，不得包含非法字符例如斜杠
	Hash       []byte      // 文件哈希，算法暂定为BLAKE2b
	UploaderId interface{} // 上传者用户ID
	CreatedAt  *gtime.Time // 文件上传时间
}
