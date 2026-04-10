package utils

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
)

// HashMagic 对设备指纹进行hash
func HashMagic(magic string) string {
	inner := sha256.Sum256([]byte(magic))
	outer := sha256.Sum256([]byte(magic + fmt.Sprintf("%x", inner)))
	return fmt.Sprintf("%x", outer)
}

// CompareMagic 校验设备指纹
func CompareMagic(magic string, hash string) bool {
	return subtle.ConstantTimeCompare([]byte(HashMagic(magic)), []byte(hash)) == 1
}
