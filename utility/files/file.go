package files

import (
	"context"
	"encoding/hex"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"golang.org/x/crypto/blake2b"
	"io"
	"log"
	"mime/multipart"
)

func HashFile(file multipart.File) []byte {
	hasher, _ := blake2b.New256(nil)
	defer file.Close()
	if _, err := io.Copy(hasher, file); err != nil {
		log.Fatal(err)
	}
	hash := hasher.Sum(nil)
	return hash
}

func HashToString(hash []byte) string {
	return hex.EncodeToString(hash[:])
}

func GetURL(fileHash []byte, fileName string) (URL *string, err error) {
	fileExtension := gstr.StrEx(fileName, ".")
	fileHashString := HashToString(fileHash)
	if fileName != "" {
		fileName = fileHashString + "." + fileExtension
	} else {
		fileName = fileHashString
	}
	ctx := context.TODO()
	uploadPath := g.Cfg().MustGet(ctx, "files.path").String()
	if uploadPath == "" {
		return nil, gerror.New("配置不存在，请配置文件地址")
	}
	*URL = "/" + uploadPath + "/" + fileName[0:2] + "/" + fileName[2:4] + "/" + fileName
	return URL, nil
}
