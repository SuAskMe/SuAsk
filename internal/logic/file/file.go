package file

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"strconv"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/model/entity"
	"suask/internal/service"
	files "suask/utility/files"
	"time"
)

type sFile struct{}

func (s *sFile) UploadFile(ctx context.Context, in model.FileUploadInput) (out *model.FileUploadOutput, err error) {
	// 获取文件保存地址
	uploadPath := g.Cfg().MustGet(ctx, "upload.path").String()
	if uploadPath == "" {
		return nil, gerror.New("配置不存在，请配置文件地址")
	}
	// 取得上传者Id
	upLoaderId := gconv.String(ctx.Value(consts.CtxId))
	if upLoaderId == "" {
		return nil, gerror.New("未找到上传者")
	}
	// 限制上传数量
	count, err := dao.Files.Ctx(ctx).
		Where(dao.Files.Columns().UploaderId, upLoaderId).
		WhereGTE(dao.Files.Columns().CreatedAt, gtime.Now().Add(-time.Minute)).Count()
	if err != nil {
		return nil, err
	}
	if count > consts.FileUploadMaxMinutes {
		return nil, gerror.New("上传频繁，一分钟只能传" + strconv.Itoa(consts.FileUploadMaxMinutes) + "次")
	}
	// 获取文件并求 hash
	file, err := in.File.Open()
	if err != nil {
		return nil, gerror.New("文件未上传")
	}
	fileHash := files.HashFile(file)
	filePath := gfile.Join(uploadPath,
		files.HashToString(fileHash)[0:2],
		files.HashToString(fileHash)[2:4],
	)
	_, err = in.File.Save(filePath, false)
	if err != nil {
		return nil, err
	}
	fileName := in.File.Filename
	// 入库
	data := do.Files{
		Name:       fileName,
		UploaderId: upLoaderId,
		Hash:       fileHash,
	}
	URL, err := files.GetURL(fileHash, fileName)
	if err != nil {
		return nil, err
	}
	id, err := dao.Files.Ctx(ctx).Data(data).OmitEmpty().InsertAndGetId()
	if err != nil {
		return nil, err
	}
	// 组合输出
	out = &model.FileUploadOutput{}
	out.Id = int(id)
	out.Name = fileName
	out.URL = URL
	return out, nil
}

func (s *sFile) Get(ctx context.Context, in model.FileGetInput) (out model.FileGetOutput, err error) {
	file := entity.Files{}
	err = dao.Files.Ctx(ctx).Where(dao.Files.Columns().Id, in.Id).Scan(&file)
	if err != nil {
		return model.FileGetOutput{}, err
	}
	URL, err := files.GetURL(file.Hash, file.Name)
	if err != nil {
		return model.FileGetOutput{}, err
	}
	out = model.FileGetOutput{}
	out.Name = file.Name
	out.URL = URL
	out.HashString = files.HashToString(file.Hash)
	out.CreatedAt = file.CreatedAt
	return out, nil
}

func init() {
	service.RegisterFile(New())
}

func New() *sFile {
	return &sFile{}
}
