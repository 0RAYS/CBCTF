package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Name    string `json:"name"`
	UserID  uint   `json:"id"`
	IsAdmin bool   `json:"admin"`
	X       string `json:"x"`
	jwt.RegisteredClaims
}

var secret = UUID()

// GenerateToken 生成token
func GenerateToken(id uint, name string, isAdmin bool, magic string) (tokenString string, err error) {
	claim := Claims{
		UserID:  id,
		Name:    name,
		IsAdmin: isAdmin,
		X:       HashMagic(magic),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString([]byte(secret))
	return tokenString, err
}

// ParseToken 解析token
func ParseToken(t string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(t, &Claims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	switch {
	case errors.Is(err, jwt.ErrTokenMalformed):
		err = errors.New("that's not even a token")
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		// Invalid signature
		err = errors.New("invalid signature")
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		// Token is either expired or not active yet
		err = errors.New("invalid signature")
	case token.Valid:
		if claims, ok := token.Claims.(*Claims); ok {
			return claims, nil
		}
	default:
		err = errors.New("couldn't handle this token")
	}
	return nil, errors.New("couldn't handle this token")
}
