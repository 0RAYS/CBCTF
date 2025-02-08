package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"context"
	"fmt"
	"time"
)

func SetEmailVerifyToken(userID uint, token string) (bool, string) {
	if !config.Env.Redis.On {
		return false, "RedisOff"
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	err := RDB.Set(ctx, fmt.Sprintf("email:verify:%d", userID), token, time.Hour*24).Err()
	if err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return false, "SetEmailVerifyTokenError"
	}
	return true, "Success"
}

func GetEmailVerifyToken(userID uint) (string, bool) {
	if !config.Env.Redis.On {
		return "RedisOff", false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := RDB.Get(ctx, fmt.Sprintf("email:verify:%d", userID)).Result()
	if err != nil {
		return "GetEmailVerifyTokenError", false
	}
	return data, true
}

func DelEmailVerifyToken(userID uint) (bool, string) {
	if !config.Env.Redis.On {
		return false, "RedisOff"
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	err := RDB.Del(ctx, fmt.Sprintf("email:verify:%d", userID)).Err()
	if err != nil {
		log.Logger.Warningf("Failed to delete email verify token: %s", err)
		return false, "DelEmailVerifyTokenError"
	}
	return true, "Success"
}
