package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"os"
)

func GetFileInfoByPath(path string) (int64, string, error) {
	var size int64
	file, err := os.Open(path)
	if err != nil {
		return 0, "", err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	hash := sha256.New()
	size, err = io.Copy(hash, file)
	if err != nil {
		return 0, "", err
	}
	return size, hex.EncodeToString(hash.Sum(nil)), nil
}

func GetFileInfoByHeader(file *multipart.FileHeader) (int64, string, error) {
	var size int64
	src, err := file.Open()
	if err != nil {
		return 0, "", err
	}
	defer func(src multipart.File) {
		_ = src.Close()
	}(src)
	sha256Sum := sha256.New()
	size, err = io.Copy(sha256Sum, src)
	if err != nil {
		return 0, "", err
	}
	return size, hex.EncodeToString(sha256Sum.Sum(nil)), nil
}
