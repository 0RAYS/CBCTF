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
	emailVerifyTokenKey     = "email:verify:token:%s"
	emailVerifyUserTokenKey = "email:verify:user:%d"
	emailVerifyTokenTTL     = 30 * time.Minute
)

// SetEmailVerifyToken 设置邮箱验证 token, 时效 30 分钟
func SetEmailVerifyToken(userID uint, token string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	userTokenKey := fmt.Sprintf(emailVerifyUserTokenKey, userID)
	oldToken, err := RDB.Get(ctx, userTokenKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Logger.Warningf("Failed to get old email verify token: %s", err)
		return model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": userTokenKey, "Error": err.Error()}}
	}

	pipe := RDB.TxPipeline()
	if oldToken != "" {
		pipe.Del(ctx, fmt.Sprintf(emailVerifyTokenKey, oldToken))
	}
	pipe.Set(ctx, fmt.Sprintf(emailVerifyTokenKey, token), strconv.FormatUint(uint64(userID), 10), emailVerifyTokenTTL)
	pipe.Set(ctx, userTokenKey, token, emailVerifyTokenTTL)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(emailVerifyTokenKey, token), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func GetEmailVerifyUserID(token string) (uint, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(emailVerifyTokenKey, token)
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, model.RetVal{Msg: i18n.Redis.NotFound, Attr: map[string]any{"Key": key}}
		}
		log.Logger.Warningf("Failed to get email verify token: %s", err)
		return 0, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	userID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	return uint(userID), model.SuccessRetVal()
}

// DelEmailVerifyToken 删除邮箱验证 token
func DelEmailVerifyToken(token string, userID uint) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Del(ctx, fmt.Sprintf(emailVerifyTokenKey, token), fmt.Sprintf(emailVerifyUserTokenKey, userID)).Err(); err != nil {
		log.Logger.Warningf("Failed to delete email verify token: %s", err)
		return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": fmt.Sprintf(emailVerifyTokenKey, token), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
