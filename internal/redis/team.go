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

func GetTeamCache(key string) (model.Team, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return model.Team{}, false
	} else if err != nil {
		return model.Team{}, false
	}
	var team model.Team
	err = msgpack.Unmarshal([]byte(data), &team)
	if err != nil {
		return model.Team{}, false
	}
	log.Logger.Debug("GetTeamCache: ", team.ID)
	return team, true
}

func GetTeamsCache(key string) ([]model.Team, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		atomic.AddInt64(&CacheMiss, 1)
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var teams []model.Team
	err = msgpack.Unmarshal([]byte(data), &teams)
	if err != nil {
		return nil, false
	}
	atomic.AddInt64(&CacheHit, 1)
	log.Logger.Debug("GetTeamsCache: ", len(teams))
	return teams, true
}

func SetTeamCache(key string, team model.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := msgpack.Marshal(team)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 10*time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetTeamCache: ", team.ID)
	return nil
}

func SetTeamsCache(key string, teams []model.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	data, err := msgpack.Marshal(teams)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, key, data, 2*time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetTeamsCache: ", len(teams))
	return nil
}

func DelTeamCache(id uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, fmt.Sprintf("team:%d:*", id), 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan teams keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				return err
			}
			log.Logger.Debug("DelTeamCache: ", key)
		}
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelTeamCache")
	return nil
}

func DelTeamsCache() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*time.Duration(config.Env.Redis.Timeout))
	defer cancel()
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, "teams:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan teams keys: %s", err)
		}

		for _, key := range keys {
			if err := RDB.Del(ctx, key).Err(); err != nil {
				return err
			}
			log.Logger.Debug("DelTeamsCache: ", key)
		}
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelTeamsCache")
	return nil
}
