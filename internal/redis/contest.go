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

func GetContestCache(key string) (model.Contest, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
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
	log.Logger.Debug("GetContestCache: ", contest.ID)
	return contest, true
}

func GetContestsCache(key string) ([]model.Contest, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
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
	log.Logger.Debug("GetContestsCache: ", len(contests))
	return contests, true
}

func SetContestCache(key string, contest model.Contest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	data, err := msgpack.Marshal(contest)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetContestCache: ", contest.ID)
	return nil
}

func SetContestsCache(key string, contests []model.Contest) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	data, err := msgpack.Marshal(contests)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 1*time.Hour).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetContestsCache: ", len(contests))
	return nil
}

func DelContestCache(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
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
			log.Logger.Debug("DelContestCache: ", key)
		}
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelContestsCache")
	return nil
}

func DelContestsCache() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, "contests:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan contest keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				return err
			}
			log.Logger.Debug("DelContestsCache: ", key)
		}
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelContestsCache")
	return nil
}
