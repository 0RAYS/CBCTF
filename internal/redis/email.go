package redis

import (
	"CBCTF/internal/log"
	"context"
	"fmt"
	"time"
)

func SetEmailVerifyToken(userID uint, token string) (bool, string) {
	ctx := context.Background()
	err := RDB.Set(ctx, fmt.Sprintf("email:verify:%d", userID), token, time.Hour*24).Err()
	if err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return false, "SetEmailVerifyTokenError"
	}
	return true, "Success"
}

func GetEmailVerifyToken(userID uint) (string, bool) {
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf("email:verify:%d", userID)).Result()
	if err != nil {
		return "GetEmailVerifyTokenError", false
	}
	return data, true
}

func DelEmailVerifyToken(userID uint) (bool, string) {
	ctx := context.Background()
	err := RDB.Del(ctx, fmt.Sprintf("email:verify:%d", userID)).Err()
	if err != nil {
		log.Logger.Warningf("Failed to delete email verify token: %s", err)
		return false, "DelEmailVerifyTokenError"
	}
	return true, "Success"
}
