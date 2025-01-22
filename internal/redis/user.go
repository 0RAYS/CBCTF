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

func GetUserCache(key string) (model.User, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
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
	log.Logger.Debug("GetUserCache: ", user.ID)
	return user, true
}

func GetUsersCache(key string) ([]model.User, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
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
	log.Logger.Debug("GetUsersCache: ", len(users))
	return users, true
}

func SetUserCache(key string, user model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	data, err := msgpack.Marshal(user)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 10*time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetUserCache: ", user.ID)
	return nil
}

func SetUsersCache(key string, users []model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	data, err := msgpack.Marshal(users)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 2*time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetUsersCache: ", len(users))
	return nil
}

func DelUserCache(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
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
	log.Logger.Debug("DelUserCache: ", id)
	return nil
}

func DelUsersCache() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, "user:*", 10).Result()
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
	log.Logger.Debug("DelUsersCache")
	return nil
}
