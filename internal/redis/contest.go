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
	ContestPattern  = "c:%d:%v" // c:<contest_id>:<preload>
	ContestsPattern = "cs:%v"   // cs:<preload>
)

func GetContestCache(id uint, preload int) (model.Contest, bool) {
	if !config.Env.Redis.On {
		return model.Contest{}, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf(ContestPattern, id, preload)).Result()
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

func GetContestsCache(preload int) ([]model.Contest, bool) {
	if !config.Env.Redis.On {
		return nil, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf(ContestsPattern, preload)).Result()
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

func SetContestCache(contest model.Contest, preload int) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(contest)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, fmt.Sprintf(ContestPattern, contest.ID, preload), data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetContestCache: ", contest.ID)
	return nil
}

func SetContestsCache(contests []model.Contest, preload int) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(contests)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, fmt.Sprintf(ContestsPattern, preload), data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetContestsCache: ", len(contests))
	return nil
}

func DelContestCache(id uint) error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(fmt.Sprintf(ContestPattern, id, "*"))
}

func DelContestsCache() error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(fmt.Sprintf(ContestsPattern, "*"))
}
