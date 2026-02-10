package redis

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"time"
)

const rateLimitKey = "rl:%s:%s"

func RateLimit(path, ip string, window time.Duration) (int64, model.RetVal) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(rateLimitKey, ip, path)
	count, err := RDB.Incr(ctx, key).Result()
	if err != nil {
		log.Logger.Warningf("Failed to increment rate limit key: %v", err)
		return 0, model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": "key", "Error": err.Error()}}
	}
	if count == 1 {
		if err = RDB.Expire(ctx, key, window).Err(); err != nil {
			log.Logger.Warningf("Failed to expire rate limit key: %v", err)
			return 0, model.RetVal{Msg: i18n.Redis.DeleteError, Attr: map[string]any{"Key": "key", "Error": err.Error()}}
		}
	}
	return count, model.SuccessRetVal()
}
