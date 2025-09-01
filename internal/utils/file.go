package utils

import (
	"CBCTF/internal/log"
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
		if err = file.Close(); err != nil {
			log.Logger.Warningf("Failed to close file: %s", err)
		}
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
		if err = src.Close(); err != nil {
			log.Logger.Warningf("Failed to close file: %s", err)
		}
	}(src)
	sha256Sum := sha256.New()
	size, err = io.Copy(sha256Sum, src)
	if err != nil {
		return 0, "", err
	}
	return size, hex.EncodeToString(sha256Sum.Sum(nil)), nil
}
