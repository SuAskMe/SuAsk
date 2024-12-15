package file

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	v1 "suask/api/file/v1"
	"suask/internal/model"
	"suask/internal/service"
)

type cFile struct {
}

var File cFile

func (c *cFile) UpdateFile(ctx context.Context, req *v1.UploadFileReq) (res *v1.UploadFileRes, err error) {
	FileInfo := model.FileUploadInput{}
	err = gconv.Struct(req, &FileInfo)
	if err != nil {
		return nil, err
	}
	out, err := service.File().UploadFile(ctx, FileInfo)
	if err != nil {
		return nil, err
	}
	return &v1.UploadFileRes{Id: out.Id, Name: out.Name, URL: out.URL}, nil
}

func (c *cFile) GetFileById(ctx context.Context, req *v1.GetFileReq) (res *v1.GetFileRes, err error) {
	FileId := model.FileGetInput{}
	err = gconv.Struct(req, &FileId)
	if err != nil {
		return nil, err
	}
	out, err := service.File().Get(ctx, FileId)
	if err != nil {
		return nil, err
	}
	return &v1.GetFileRes{
		URL:        out.URL,
		Name:       out.Name,
		HashString: out.HashString,
		UploaderId: out.UploaderId,
		CreatedAt:  out.CreatedAt,
	}, nil
}
