package redis

import (
	"context"
	"fmt"
	"time"
)

const rateLimitKey = "rl:%s:%s"

func RateLimit(path, ip string, window time.Duration) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(rateLimitKey, ip, path)
	count, err := RDB.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if count == 1 {
		if err = RDB.Expire(ctx, key, window).Err(); err != nil {
			return 0, err
		}
	}
	return count, nil
}
