package netutil

import (
	"errors"
	"mime/multipart"
	"os"
)

// Size 获取文件大小的接口
type Size interface {
	Size() int64
}

// Stat 获取文件信息的接口
type Stat interface {
	Stat() (os.FileInfo, error)
}

// GetUploadFileSize 返回HTTP upload文件的大小
func GetUploadFileSize(upfile multipart.File) (int64, error) {
	if statInterface, ok := upfile.(Stat); ok {
		fileInfo, _ := statInterface.Stat()
		return fileInfo.Size(), nil
	}
	if sizeInterface, ok := upfile.(Size); ok {
		fsize := sizeInterface.Size()
		return fsize, nil
	}
	return 0, errors.New("not found stat and size interface")
}
