package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

const (
	emailVerifyTokenKeyTmpl       = "email:verify:token:%s"
	emailVerifyUserTokenKeyTmpl   = "email:verify:user:%d"
	emailVerifyTokenTTL           = 30 * time.Minute
	passwordResetTokenKeyTmpl     = "email:password:reset:token:%s"
	passwordResetUserTokenKeyTmpl = "email:password:reset:user:%d"
	passwordResetTokenTTL         = 30 * time.Minute
)

// tokenStore binds the keys and the lifetime of a one-time token that is
// stored both by token and by user.
type tokenStore struct {
	label    string
	tokenKey string
	userKey  string
	ttl      time.Duration
}

func emailVerifyTokens() tokenStore {
	return tokenStore{"email verify", emailVerifyTokenKeyTmpl, emailVerifyUserTokenKeyTmpl, emailVerifyTokenTTL}
}

func passwordResetTokens() tokenStore {
	return tokenStore{"password reset", passwordResetTokenKeyTmpl, passwordResetUserTokenKeyTmpl, passwordResetTokenTTL}
}

// setToken stores the token and drops the token the user held before.
func (t tokenStore) setToken(userID uint, token string) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	userTokenKey := fmt.Sprintf(t.userKey, userID)
	oldToken, err := RDB.Get(ctx, userTokenKey).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Logger.Warningf("Failed to get old %s token: %s", t.label, err)
		return model.FailRetVal(i18n.Redis.GetError, "Cache", userTokenKey, err)
	}

	pipe := RDB.TxPipeline()
	if oldToken != "" {
		pipe.Del(ctx, fmt.Sprintf(t.tokenKey, oldToken))
	}
	pipe.Set(ctx, fmt.Sprintf(t.tokenKey, token), strconv.FormatUint(uint64(userID), 10), t.ttl)
	pipe.Set(ctx, userTokenKey, token, t.ttl)
	if _, err := pipe.Exec(ctx); err != nil {
		log.Logger.Warningf("Failed to set %s token: %s", t.label, err)
		return model.FailRetVal(i18n.Redis.SetError, "Cache", fmt.Sprintf(t.tokenKey, token), err)
	}
	return model.SuccessRetVal()
}

// getUserID resolves the user a token was issued for.
func (t tokenStore) getUserID(token string) (uint, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(t.tokenKey, token)
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, model.FailRetVal(i18n.Redis.NotFound, "Cache", key, nil)
		}
		log.Logger.Warningf("Failed to get %s token: %s", t.label, err)
		return 0, model.FailRetVal(i18n.Redis.GetError, "Cache", key, err)
	}
	userID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, model.FailRetVal(i18n.Redis.GetError, "Cache", key, err)
	}
	return uint(userID), model.SuccessRetVal()
}

// delToken removes the token and the user index of it.
func (t tokenStore) delToken(token string, userID uint) model.RetVal {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	if err := RDB.Del(ctx, fmt.Sprintf(t.tokenKey, token), fmt.Sprintf(t.userKey, userID)).Err(); err != nil {
		log.Logger.Warningf("Failed to delete %s token: %s", t.label, err)
		return model.FailRetVal(i18n.Redis.DeleteError, "Cache", fmt.Sprintf(t.tokenKey, token), err)
	}
	return model.SuccessRetVal()
}

// SetEmailVerifyToken 设置邮箱验证 token, 时效 30 分钟
func SetEmailVerifyToken(userID uint, token string) model.RetVal {
	return emailVerifyTokens().setToken(userID, token)
}

func GetEmailVerifyUserID(token string) (uint, model.RetVal) {
	return emailVerifyTokens().getUserID(token)
}

// DelEmailVerifyToken 删除邮箱验证 token
func DelEmailVerifyToken(token string, userID uint) model.RetVal {
	return emailVerifyTokens().delToken(token, userID)
}

// SetPasswordResetToken 设置密码重置 token, 时效 30 分钟
func SetPasswordResetToken(userID uint, token string) model.RetVal {
	return passwordResetTokens().setToken(userID, token)
}

func GetPasswordResetUserID(token string) (uint, model.RetVal) {
	return passwordResetTokens().getUserID(token)
}

// DelPasswordResetToken 删除密码重置 token
func DelPasswordResetToken(token string, userID uint) model.RetVal {
	return passwordResetTokens().delToken(token, userID)
}
