package utils

import (
	"io"
	"mime/multipart"
	"os"
)

// 读取文件为[]byte类型
func ReadFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()
	// 读取文件内容
	return io.ReadAll(file)
}

// 保存一个文件
func SaveAFile(savePath string, file multipart.File) error {
	out, err := os.Create(savePath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	return err
}
