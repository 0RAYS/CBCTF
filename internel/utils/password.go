package utils

import (
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

const (
	VeryWeak = 0
	Weak     = 1
	Medium   = 2
	Strong   = 3
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
	var length, t, upper, lower, digit, special int

	length = len(password)
	if length < 8 {
		return VeryWeak
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upper++
		case unicode.IsLower(char):
			lower++
		case unicode.IsDigit(char):
			digit++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special++
		}
	}
	if upper > 0 {
		t++
	}
	if lower > 0 {
		t++
	}
	if digit > 0 {
		t++
	}
	if special > 0 {
		t++
	}
	switch {
	case t == 4:
		return Strong
	case t == 3:
		return Medium
	case t == 2:
		return Weak
	default:
		return VeryWeak
	}
}
