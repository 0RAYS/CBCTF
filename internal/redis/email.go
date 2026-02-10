package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"time"
)

const emailVerifyTokenKey = "email:%d"

// SetEmailVerifyToken 设置邮箱验证 token, 时效一天
func SetEmailVerifyToken(userID uint, token string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Set(ctx, fmt.Sprintf(emailVerifyTokenKey, userID), token, 30*time.Minute).Err(); err != nil {
		log.Logger.Warningf("Failed to set email verify token: %s", err)
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": fmt.Sprintf(emailVerifyTokenKey, userID), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

// GetEmailVerifyToken 获取邮箱验证 token
func GetEmailVerifyToken(userID uint) (string, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	token, err := RDB.Get(ctx, fmt.Sprintf(emailVerifyTokenKey, userID)).Result()
	if err != nil {
		log.Logger.Warningf("Failed to get email verify token: %s", err)
		return token, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": fmt.Sprintf(emailVerifyTokenKey, userID), "Error": err.Error()}}
	}
	return token, model.SuccessRetVal()
}

// DelEmailVerifyToken 删除邮箱验证 token
func DelEmailVerifyToken(userID uint) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Del(ctx, fmt.Sprintf(emailVerifyTokenKey, userID)).Err(); err != nil {
		log.Logger.Warningf("Failed to delete email verify token: %s", err)
		return model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": fmt.Sprintf(emailVerifyTokenKey, userID), "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
