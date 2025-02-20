package redis

import (
	"context"
)

// DeleteKeysByPattern 删除匹配模式的缓存键
func DeleteKeysByPattern(pattern string) error {
	ctx := context.Background()
	var cursor uint64
	for {
		keys, nextCursor, err := RDB.Scan(ctx, cursor, pattern, 1000).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := RDB.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}
