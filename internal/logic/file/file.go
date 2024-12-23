package file

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
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

type File struct {
	id   int
	name string
	URL  string
}

func uploadFile(upLoadFile *ghttp.UploadFile, upLoaderId int) (fileInfo *File, err error) {
	file, err := upLoadFile.Open()
	if err != nil {
		return nil, gerror.New("文件未上传")
	}
	uploadPath := g.Cfg().MustGet(context.TODO(), "upload.path").String()
	if uploadPath == "" {
		return nil, gerror.New("配置不存在，请配置文件地址")
	}
	fileHash := files.HashFile(file)
	filePath := gfile.Join(uploadPath,
		files.HashToString(fileHash)[0:2],
		files.HashToString(fileHash)[2:4],
	)
	oldName := upLoadFile.Filename
	newName, err := files.RenameFiles(fileHash, upLoadFile.Filename)
	if err != nil {
		return nil, err
	}
	upLoadFile.Filename = newName
	_, err = upLoadFile.Save(filePath, false)
	if err != nil {
		return nil, err
	}
	fileName := oldName
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
	id, err := dao.Files.Ctx(context.TODO()).Data(data).OmitEmpty().InsertAndGetId()
	if err != nil {
		return nil, err
	}
	fileInfo = &File{
		id:   int(id),
		name: fileName,
		URL:  URL,
	}
	return fileInfo, nil
}

func (s *sFile) UploadFile(ctx context.Context, in model.FileUploadInput) (out *model.FileUploadOutput, err error) {
	// 取得上传者Id
	upLoaderId := gconv.Int(ctx.Value(consts.CtxId))
	if upLoaderId == 0 {
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
	fileInfo, err := uploadFile(in.File, upLoaderId)
	if err != nil {
		return nil, err
	}
	// 组合输出
	out = &model.FileUploadOutput{
		Id:   fileInfo.id,
		Name: fileInfo.name,
		URL:  fileInfo.URL,
	}
	return out, nil
}

func (s *sFile) UploadFileList(_ context.Context, in model.FileListAddInput) (out model.FileListAddOutput, err error) {
	fileCount := len(in.FileList)
	out = model.FileListAddOutput{
		IdList: make([]int, fileCount),
	}
	for fileIndex, file := range in.FileList {
		fileInfo, err := uploadFile(file, in.UploaderId)
		if err != nil {
			return model.FileListAddOutput{}, err
		}
		out.IdList[fileIndex] = fileInfo.id
		fmt.Println(out)
	}
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

func (s *sFile) GetList(ctx context.Context, in model.FileListGetInput) (out model.FileListGetOutput, err error) {
	var fileList []entity.Files
	var count int
	err = dao.Files.Ctx(ctx).WhereIn(dao.Files.Columns().Id, in.IdList).Order(dao.Files.Columns().Id).ScanAndCount(&fileList, &count, false)
	if err != nil {
		return model.FileListGetOutput{}, err
	}
	fmt.Println(count)
	out = model.FileListGetOutput{
		FileId:     make([]int, count),
		Name:       make([]string, count),
		URL:        make([]string, count),
		UploaderId: make([]int, count),
		CreatedAt:  make([]*gtime.Time, count),
	}
	for index, file := range fileList {
		out.Name[index] = file.Name
		URL, err := files.GetURL(file.Hash, file.Name)
		if err != nil {
			return model.FileListGetOutput{}, err
		}
		out.FileId[index] = file.Id
		out.URL[index] = URL
		out.UploaderId[index] = file.UploaderId
		out.CreatedAt[index] = file.CreatedAt
	}
	return out, nil
}

func init() {
	service.RegisterFile(New())
}

func New() *sFile {
	return &sFile{}
}
