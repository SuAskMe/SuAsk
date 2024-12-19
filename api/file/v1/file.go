package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

type UploadFileReq struct {
	g.Meta `path:"/files" method:"post" tags:"Files" summary:"上传图片"`
	File   *ghttp.UploadFile `json:"file" v:"required" dc:"要上传的图片"`
}

type UploadFileRes struct {
	Id   int    `json:"id" dc:"文件id"`
	Name string `json:"name" dc:"文件名，包含文件后缀"`
	URL  string `json:"url" dc:"返回给前端的url"`
}

type GetFileReq struct {
	g.Meta `path:"/files" method:"get" tags:"Files" summary:"通过文件Id获取文件"`
	Id     int `json:"id" dc:"文件id"`
}
type GetFileRes struct {
	URL        string      `json:"url" dc:"文件的URl"`
	Name       string      `json:"name" dc:"文件原始名称"`
	HashString string      `json:"hash" dc:"文件的Hash值"`
	UploaderId int         `json:"uploader_id" dc:"文件上传者的id"`
	CreatedAt  *gtime.Time `json:"created_at" dc:"文件的上传时间"`
}

type GetFileListReq struct {
	g.Meta `path:"/file-list" method:"get" tags:"Files" summary:"通过文件IdList获取文件"`
	Id     []int `json:"id"`
}

type GetFileListRes struct {
	URL        []string      `json:"url" dc:"文件的URl"`
	Name       []string      `json:"name" dc:"文件原始名称"`
	UploaderId []int         `json:"uploader_id" dc:"文件上传者的id"`
	CreatedAt  []*gtime.Time `json:"created_at" dc:"文件的上传时间"`
}
