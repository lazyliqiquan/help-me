package utils

import (
	"io"
	"mime/multipart"
	"os"
)

// ReadFile 读取文件为[]byte类型
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	// 读取文件内容
	return io.ReadAll(file)
}

// SaveAFile 保存一个文件
func SaveAFile(savePath string, file multipart.File) error {
	// 这里关闭的文件应该不总是最后一个吧(如果程序内存溢出可以考虑文件是否及时关闭)
	defer func() {
		err := file.Close()
		if err != nil {
			Logger.Errorln(err)
		}
	}()
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}
