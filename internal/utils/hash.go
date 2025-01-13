package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 生成密码hash
func HashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}

// CompareHashAndPassword 验证密码
func CompareHashAndPassword(hash string, password string) bool {
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return false
	}
	return true
}
