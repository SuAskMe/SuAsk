// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// Themes is the golang structure for table themes.
type Themes struct {
	Id               int `json:"id"               orm:"id"                 description:"主题ID"`     // 主题ID
	BackgroundFileId int `json:"backgroundFileId" orm:"background_file_id" description:"背景图片文件ID"` // 背景图片文件ID
}
