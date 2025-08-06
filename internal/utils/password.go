package utils

import (
	"unicode"

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

// CheckPassword 检查密码强度
func CheckPassword(password string) uint {
	var t, u, l, d, s uint

	if len(password) < 8 {
		return t
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			u++
		case unicode.IsLower(char):
			l++
		case unicode.IsDigit(char):
			d++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			s++
		}
	}
	if u > 0 {
		t++
	}
	if l > 0 {
		t++
	}
	if d > 0 {
		t++
	}
	if s > 0 {
		t++
	}
	return t
}
