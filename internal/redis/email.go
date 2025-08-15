package redis

import (
	"context"
	"fmt"
	"time"
)

const emailVerifyTokenKey = "email:%d"

// SetEmailVerifyToken 设置邮箱验证 token, 时效一天
func SetEmailVerifyToken(userID uint, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Set(ctx, fmt.Sprintf(emailVerifyTokenKey, userID), token, 30*time.Minute).Err()
}

// GetEmailVerifyToken 获取邮箱验证 token
func GetEmailVerifyToken(userID uint) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Get(ctx, fmt.Sprintf(emailVerifyTokenKey, userID)).Result()
}

// DelEmailVerifyToken 删除邮箱验证 token
func DelEmailVerifyToken(userID uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return RDB.Del(ctx, fmt.Sprintf(emailVerifyTokenKey, userID)).Err()
}
