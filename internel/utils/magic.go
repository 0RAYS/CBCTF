package utils

import (
	"crypto/md5"
	"fmt"
)

// HashMagic 对设备指纹进行hash
func HashMagic(magic string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(magic+fmt.Sprintf("%x", md5.Sum([]byte(magic+secret))))))
}

// CompareMagic 校验设备指纹
func CompareMagic(magic string, hash string) bool {
	return HashMagic(magic) == hash
}
