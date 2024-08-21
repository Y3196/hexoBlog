package utils

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"path/filepath"
)

// GetMd5 获取文件的 MD5 值
func GetMd5(inputStream io.Reader) (string, error) {
	hash := md5.New()
	if _, err := io.Copy(hash, inputStream); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetExtName 获取文件扩展名
func GetExtName(fileName string) string {
	return filepath.Ext(fileName)
}
