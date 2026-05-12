package file

import (
	"context"
	"suask/internal/consts"
	"suask/internal/dao"
	"suask/internal/model"
	"suask/internal/model/do"
	"suask/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/ghttp"
)

// UploadAvatarInput 头像上传输入参数
type UploadAvatarInput struct {
	UserId int               // 目标用户 ID
	File   *ghttp.UploadFile // 上传的头像文件
}

// UploadAvatarOutput 头像上传输出
type UploadAvatarOutput struct {
	UserId    int
	AvatarURL string
}

// UploadAvatar 通用头像上传逻辑，供普通用户和管理员共用。
// 1. 上传文件到存储
// 2. 更新 users.avatar_file_id
// 3. 返回新头像 URL
func UploadAvatar(ctx context.Context, in UploadAvatarInput) (out *UploadAvatarOutput, err error) {
	if in.File == nil {
		return nil, gerror.New("请上传头像文件")
	}

	// 上传文件
	fileData, err := service.File().UploadFile(ctx, model.FileUploadInput{File: in.File})
	if err != nil {
		return nil, gerror.Wrap(err, "头像上传失败")
	}

	// 更新用户头像 ID
	_, err = dao.Users.Ctx(ctx).
		Where(dao.Users.Columns().Id, in.UserId).
		Data(do.Users{AvatarFileId: fileData.Id}).
		Update()
	if err != nil {
		return nil, gerror.New(consts.ErrInternal)
	}

	// 获取新头像 URL
	newFile, err := service.File().Get(ctx, model.FileGetInput{Id: fileData.Id})
	avatarURL := consts.DefaultAvatarURL
	if err == nil {
		avatarURL = newFile.URL
	}

	out = &UploadAvatarOutput{
		UserId:    in.UserId,
		AvatarURL: avatarURL,
	}
	return
}
