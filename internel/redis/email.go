package redis

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/log"
	"context"
	"fmt"
	"time"
)

// SetEmailVerifyToken 设置邮箱验证 token, 时效一天
func SetEmailVerifyToken(userID uint, token string) (bool, string) {
	ctx := context.Background()
	err := RDB.Set(ctx, fmt.Sprintf("email:%d", userID), token, time.Hour*24).Err()
	if err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return false, i18n.SetEmailVerifyTokenError
	}
	return true, i18n.Success
}

// GetEmailVerifyToken 获取邮箱验证 token
func GetEmailVerifyToken(userID uint) (string, bool) {
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf("email:%d", userID)).Result()
	if err != nil {
		return i18n.GetEmailVerifyTokenError, false
	}
	return data, true
}

// DelEmailVerifyToken 删除邮箱验证 token
func DelEmailVerifyToken(userID uint) (bool, string) {
	ctx := context.Background()
	err := RDB.Del(ctx, fmt.Sprintf("email:%d", userID)).Err()
	if err != nil {
		log.Logger.Warningf("Failed to delete email verify token: %s", err)
		return false, i18n.DelEmailVerifyTokenError
	}
	return true, i18n.Success
}
