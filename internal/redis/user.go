package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"sync/atomic"
	"time"
)

const (
	UserPattern  = "u:%d:%v" // user:<user_id>:<preload>
	UsersPattern = "us:%v"   // users:<preload>
)

func GetUserCache(id uint, preload int) (model.User, bool) {
	if !config.Env.Redis.On {
		return model.User{}, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf(UserPattern, id, preload)).Result()
	if errors.Is(err, redis.Nil) {
		atomic.AddInt64(&CacheMiss, 1)
		return model.User{}, false
	} else if err != nil {
		return model.User{}, false
	}
	var user model.User
	err = msgpack.Unmarshal([]byte(data), &user)
	if err != nil {
		return model.User{}, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetUserCache: ", user.ID)
	return user, true
}

func GetUsersCache(preload int) ([]model.User, bool) {
	if !config.Env.Redis.On {
		return nil, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf(UsersPattern, preload)).Result()
	if errors.Is(err, redis.Nil) {
		atomic.AddInt64(&CacheMiss, 1)
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var users []model.User
	err = msgpack.Unmarshal([]byte(data), &users)
	if err != nil {
		return nil, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetUsersCache: ", len(users))
	return users, true
}

func SetUserCache(user model.User, preload int) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(user)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, fmt.Sprintf(UserPattern, user.ID, preload), data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetUserCache: ", user.ID)
	return nil
}

func SetUsersCache(users []model.User, preload int) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(users)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, fmt.Sprintf(UsersPattern, preload), data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetUsersCache: ", len(users))
	return nil
}

func DelUserCache(id uint) error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(fmt.Sprintf(UserPattern, id, "*"))
}

func DelUsersCache() error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(fmt.Sprintf(UsersPattern, "*"))
}
