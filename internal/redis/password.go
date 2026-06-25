package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	passwordResetTokenKey = "password:reset:%d"
	passwordResetTokenTTL = 30 * time.Minute
)

// SetPasswordResetToken 设置密码重置 token, 时效 30 分钟
func SetPasswordResetToken(userID uint, token string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Set(ctx, fmt.Sprintf(passwordResetTokenKey, userID), token, passwordResetTokenTTL).Err(); err != nil {
		log.Logger.Warningf("Failed to set password reset token: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(passwordResetTokenKey, userID), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// GetPasswordResetToken 获取密码重置 token
func GetPasswordResetToken(userID uint) (string, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	token, err := RDB.Get(ctx, fmt.Sprintf(passwordResetTokenKey, userID)).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", model.RetVal{Msg: i18n.Redis.NotFound, Attr: map[string]any{"Key": fmt.Sprintf(passwordResetTokenKey, userID)}}
		}
		log.Logger.Warningf("Failed to get password reset token: %s", err)
		return "", model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": fmt.Sprintf(passwordResetTokenKey, userID), "Error": err.Error()}}
	}
	return token, model.SuccessRetVal()
}

// DelPasswordResetToken 删除密码重置 token
func DelPasswordResetToken(userID uint) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Del(ctx, fmt.Sprintf(passwordResetTokenKey, userID)).Err(); err != nil {
		log.Logger.Warningf("Failed to delete password reset token: %s", err)
		return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": fmt.Sprintf(passwordResetTokenKey, userID), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
