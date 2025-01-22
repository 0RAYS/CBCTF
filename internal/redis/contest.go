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

func GetContestCache(ctx context.Context, key string) (model.Contest, bool) {
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Contest{}, false
	} else if err != nil {
		return model.Contest{}, false
	}
	var contest model.Contest
	err = msgpack.Unmarshal([]byte(data), &contest)
	if err != nil {
		return model.Contest{}, false
	}
	return contest, true
}

func GetContestsCache(ctx context.Context, key string) ([]model.Contest, bool) {
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var contests []model.Contest
	err = msgpack.Unmarshal([]byte(data), &contests)
	if err != nil {
		return nil, false
	}
	return contests, true
}

func SetContestCache(ctx context.Context, key string, contest model.Contest) error {
	data, err := msgpack.Marshal(contest)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func SetContestsCache(ctx context.Context, key string, contests []model.Contest) error {
	data, err := msgpack.Marshal(contests)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func DelContestCache(ctx context.Context, id uint) error {
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, fmt.Sprintf("contest:%d:*", id), 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan contest keys: %s", err)
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

func DelContestsCache(ctx context.Context) error {
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, "contest:list:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan contest keys: %s", err)
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
