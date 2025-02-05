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

func GetContestCache(key string) (model.Contest, bool) {
	if !config.Env.Redis.On {
		return model.Contest{}, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		atomic.AddInt64(&CacheMiss, 1)
		return model.Contest{}, false
	} else if err != nil {
		return model.Contest{}, false
	}
	var contest model.Contest
	err = msgpack.Unmarshal([]byte(data), &contest)
	if err != nil {
		return model.Contest{}, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetContestCache: ", contest.ID)
	return contest, true
}

func GetContestsCache(key string) ([]model.Contest, bool) {
	if !config.Env.Redis.On {
		return nil, false
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		atomic.AddInt64(&CacheMiss, 1)
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var contests []model.Contest
	err = msgpack.Unmarshal([]byte(data), &contests)
	if err != nil {
		return nil, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetContestsCache: ", len(contests))
	return contests, true
}

func SetContestCache(key string, contest model.Contest) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := msgpack.Marshal(contest)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetContestCache: ", contest.ID)
	return nil
}

func SetContestsCache(key string, contests []model.Contest) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := msgpack.Marshal(contests)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetContestsCache: ", len(contests))
	return nil
}

func DelContestCache(id uint) error {
	if !config.Env.Redis.On {
		return nil
	}
	var cursor uint64
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
		keys, cursor, err := RDB.Scan(ctx, cursor, fmt.Sprintf("contest:%d:*", id), 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan contest keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				cancel()
				return err
			}
			log.Logger.Debug("DelContestCache: ", key)
		}
		cancel()
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelContestsCache")
	return nil
}

func DelContestsCache() error {
	if !config.Env.Redis.On {
		return nil
	}
	var cursor uint64
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
		keys, cursor, err := RDB.Scan(ctx, cursor, "contests:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan contest keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				cancel()
				return err
			}
			log.Logger.Debug("DelContestsCache: ", key)
		}
		cancel()
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelContestsCache")
	return nil
}
