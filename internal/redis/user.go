package redis

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"time"
)

func GetUserCache(ctx context.Context, key string) (model.User, bool) {
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.User{}, false
	} else if err != nil {
		return model.User{}, false
	}
	var user model.User
	err = msgpack.Unmarshal([]byte(data), &user)
	if err != nil {
		return model.User{}, false
	}
	return user, true
}

func GetUsersCache(ctx context.Context, key string) ([]model.User, bool) {
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var users []model.User
	err = msgpack.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil, false
	}
	return users, true
}

func SetUserCache(ctx context.Context, key string, user model.User) error {
	data, err := msgpack.Marshal(user)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 10*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}

func SetUsersCache(ctx context.Context, key string, users []model.User) error {
	data, err := msgpack.Marshal(users)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 2*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}

func DelUserCache(ctx context.Context, id uint) error {
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, fmt.Sprintf("user:%d:*", id), 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan user keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}

func DelUsersCache(ctx context.Context) error {
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, "user:list:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan users keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	return nil
}
