package utils

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Claims struct {
	Name   string `json:"name"`
	UserID uint   `json:"id"`
	Type   string `json:"type"`
	X      string `json:"x"`
	jwt.RegisteredClaims
}

// var secret = uuid.New().String()
var secret = "0RAYS-JBNRZ"

// Generate 生成token
func Generate(id uint, name string, t string, magic string) (tokenString string, err error) {
	claim := Claims{
		UserID: id,
		Name:   name,
		Type:   t,
		X:      HashMagic(magic),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokenString, err = token.SignedString([]byte(secret))
	return tokenString, err
}

// Parse 解析token
func Parse(t string) (*Claims, error) {
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
