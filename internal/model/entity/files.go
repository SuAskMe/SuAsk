// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Files is the golang structure for table files.
type Files struct {
	Id         int    `json:"id"         orm:"id"          description:"文件ID"`              // 文件ID
	Name       string `json:"name"       orm:"name"        description:"文件名，不得包含非法字符例如斜杠"`  // 文件名，不得包含非法字符例如斜杠
	Hash       []byte `json:"hash"       orm:"hash"        description:"文件哈希，算法暂定为BLAKE2b"` // 文件哈希，算法暂定为BLAKE2b
	UploaderId int    `json:"uploaderId" orm:"uploader_id" description:"上传者用户ID"`           // 上传者用户ID
}
