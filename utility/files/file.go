package files

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/crypto/blake2b"
)

// HashFile 读取 multipart.File 的全部内容并返回 BLAKE2b-256 摘要。
// 之前这里 io.Copy 失败时会 log.Fatal 掉整个进程；现在改为返回 error，
// 把是否继续由调用方决定。
func HashFile(file multipart.File) ([]byte, error) {
	hasher, _ := blake2b.New256(nil)
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func HashToString(hash []byte) string {
	return hex.EncodeToString(hash[:])
}

func RenameFiles(fileHash []byte, fileName string) (newName string, err error) {
	fileExtension := gstr.StrEx(fileName, ".")
	fileHashString := HashToString(fileHash)
	if fileExtension != "" {
		fileName = fileHashString + "." + fileExtension
	} else {
		fileName = fileHashString
	}
	newName = fileName
	return
}

func GetURL(fileHash []byte, fileName string) (URL string, err error) {
	fileExtension := gstr.StrEx(fileName, ".")
	fileHashString := HashToString(fileHash)
	if fileExtension != "" {
		fileName = fileHashString + "." + fileExtension
	} else {
		fileName = fileHashString
	}
	ctx := context.TODO()
	uploadPath := g.Cfg().MustGet(ctx, "upload.path").String()
	if uploadPath == "" {
		return "", gerror.New("配置不存在，请配置文件地址")
	}
	var serverPrefix = g.Cfg().MustGet(context.TODO(), "upload.prefix").String()
	URL = serverPrefix + "/" + uploadPath + "/" + fileName[0:2] + "/" + fileName[2:4] + "/" + fileName
	return URL, nil
}

// --- 上传安全校验 -------------------------------------------------------------

// 默认的允许上传类型白名单。只覆盖当前业务需要的图片格式。
// 如需扩展（比如音频/视频），在 config.yaml 里加 upload.allowed_mimes 列表，
// 而不是改这里的默认值。
var defaultAllowedMimes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"image/webp",
}

// 默认最大 32 MiB；可用 config.yaml 里的 upload.max_bytes 覆盖。
// 注意：GoFrame 的 server.clientMaxBodySize 默认是 8 MiB，请求体在进入应用前就会被截掉，
// 所以如果这里调大，也要同步把 server.clientMaxBodySize 设到不小于这个值，
// 否则请求在 HTTP 层就被拒，业务层的上限再大都拦不到。
const defaultMaxUploadBytes = 32 * 1024 * 1024

// UploadPolicy 约束一次上传请求：文件大小 + MIME 类型白名单。
type UploadPolicy struct {
	MaxBytes     int64
	AllowedMimes []string
}

// LoadUploadPolicy 按优先级 config > 默认值 组装策略。读配置用 TryGet，
// 避免配置缺失时整个流程崩掉。
func LoadUploadPolicy(ctx context.Context) UploadPolicy {
	p := UploadPolicy{
		MaxBytes:     defaultMaxUploadBytes,
		AllowedMimes: append([]string(nil), defaultAllowedMimes...),
	}
	cfg := g.Cfg()
	if v, err := cfg.Get(ctx, "upload.max_bytes"); err == nil && !v.IsNil() && v.Int64() > 0 {
		p.MaxBytes = v.Int64()
	}
	if v, err := cfg.Get(ctx, "upload.allowed_mimes"); err == nil && !v.IsNil() {
		if list := v.Strings(); len(list) > 0 {
			p.AllowedMimes = list
		}
	}
	return p
}

// ValidateUpload 对上传的单个文件做"门禁":
//   - 大小不超过策略上限
//   - 通过魔术字节嗅探出的 MIME 在白名单里
//
// 注意：这里只读前 512 字节（http.DetectContentType 的固定窗口），不会把
// 整个文件读进内存；调用方随后仍需要通过 Open() 重新读取整份数据。
func ValidateUpload(fh *ghttp.UploadFile, policy UploadPolicy) error {
	if fh == nil || fh.FileHeader == nil {
		return gerror.New("文件未上传")
	}
	if fh.Size > policy.MaxBytes {
		return fmt.Errorf("文件过大（上限 %d 字节）", policy.MaxBytes)
	}

	f, err := fh.Open()
	if err != nil {
		return gerror.New("无法打开上传的文件")
	}
	defer f.Close()

	head := make([]byte, 512)
	n, err := io.ReadFull(f, head)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		return err
	}
	head = head[:n]

	mime := http.DetectContentType(head)
	// DetectContentType 总会返回带可选参数的 MIME（如 "text/plain; charset=utf-8"），
	// 拿分号前段做精确匹配。
	if semi := strings.IndexByte(mime, ';'); semi >= 0 {
		mime = strings.TrimSpace(mime[:semi])
	}
	if !mimeAllowed(mime, policy.AllowedMimes) {
		return fmt.Errorf("不支持的文件类型：%s", mime)
	}

	// 再做一次扩展名自洽检查，避免 MIME 通过但扩展名是 .exe / .html 一类。
	if err := extensionLooksSafe(fh.Filename, mime); err != nil {
		return err
	}
	return nil
}

func mimeAllowed(mime string, whitelist []string) bool {
	for _, m := range whitelist {
		if strings.EqualFold(mime, m) {
			return true
		}
	}
	return false
}

// 扩展名白名单：只和 image/* 对齐。未知扩展名不阻塞（空扩展允许），
// 但 .exe / .bat / .html / .js / .php 这类一定拒绝。
var dangerousExts = map[string]struct{}{
	".exe": {}, ".bat": {}, ".cmd": {}, ".com": {}, ".msi": {},
	".sh": {}, ".ps1": {},
	".html": {}, ".htm": {}, ".js": {}, ".mjs": {}, ".svg": {},
	".php": {}, ".asp": {}, ".aspx": {}, ".jsp": {},
}

func extensionLooksSafe(name, mime string) error {
	ext := strings.ToLower(filepath.Ext(name))
	if _, bad := dangerousExts[ext]; bad {
		return fmt.Errorf("禁止的文件扩展名：%s", ext)
	}
	// 只有 image/* 时，如果扩展名存在，要求是已知图片扩展名
	if strings.HasPrefix(mime, "image/") && ext != "" {
		switch ext {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp":
			return nil
		default:
			return fmt.Errorf("图片文件扩展名异常：%s", ext)
		}
	}
	return nil
}
