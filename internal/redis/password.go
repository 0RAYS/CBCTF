package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	passwordResetTokenKey     = "password:reset:token:%s"
	passwordResetUserTokenKey = "password:reset:user:%d"
	passwordResetTokenTTL     = 30 * time.Minute
)

// SetPasswordResetToken 设置密码重置 token, 时效 30 分钟
func SetPasswordResetToken(userID uint, token string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	userKey := fmt.Sprintf(passwordResetUserTokenKey, userID)
	oldToken, err := RDB.Get(ctx, userKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Logger.Warningf("Failed to get old password reset token: %s", err)
		return model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": userKey, "Error": err.Error()}}
	}

	pipe := RDB.TxPipeline()
	if oldToken != "" {
		pipe.Del(ctx, fmt.Sprintf(passwordResetTokenKey, oldToken))
	}
	pipe.Set(ctx, fmt.Sprintf(passwordResetTokenKey, token), strconv.FormatUint(uint64(userID), 10), passwordResetTokenTTL)
	pipe.Set(ctx, userKey, token, passwordResetTokenTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Logger.Warningf("Failed to set password reset token: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(passwordResetTokenKey, token), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetPasswordResetUserID(token string) (uint, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(passwordResetTokenKey, token)
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, model.RetVal{Msg: i18n.Redis.NotFound, Attr: map[string]any{"Key": key}}
		}
		log.Logger.Warningf("Failed to get password reset token: %s", err)
		return 0, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	userID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	return uint(userID), model.SuccessRetVal()
}

// DelPasswordResetToken 删除密码重置 token
func DelPasswordResetToken(token string, userID uint) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Del(ctx, fmt.Sprintf(passwordResetTokenKey, token), fmt.Sprintf(passwordResetUserTokenKey, userID)).Err(); err != nil {
		log.Logger.Warningf("Failed to delete password reset token: %s", err)
		return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": fmt.Sprintf(passwordResetTokenKey, token), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
