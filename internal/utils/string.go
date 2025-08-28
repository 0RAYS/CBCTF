package utils

import (
	"crypto/rand"
	"encoding/hex"
	"unicode"

	"github.com/google/uuid"
)

// UUID 生成随机uuid
func UUID() string {
	return uuid.New().String()
}

// RandStr 生成随机字符串
func RandStr(n int) string {
	result := make([]byte, n)
	_, err := rand.Read(result)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(result)[:n]
}

func ToTitle(s string) string {
	if len(s) == 0 {
		return ""
	}
	str := []rune(s)
	str[0] = unicode.ToTitle(str[0])
	return string(str)
}
