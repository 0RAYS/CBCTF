package redis

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"time"
)

func GetTeamCache(key string) (model.Team, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	data, err := RDB.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	} else if err != nil {
		return nil, false
	}
	var teams []model.Team
	err = msgpack.Unmarshal([]byte(data), &teams)
	if err != nil {
		return nil, false
	}
	log.Logger.Debug("GetTeamsCache: ", len(teams))
	return teams, true
}

func SetTeamCache(key string, team model.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
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

func DelTeamsCache() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	var cursor uint64
	for {
		keys, cursor, err := RDB.Scan(ctx, cursor, "team:*", 10).Result()
		if err != nil {
			log.Logger.Warningf("Failed to scan teams keys: %s", err)
		}

		for _, key := range keys {
			log.Logger.Debug("DelTeamsCache: ", key)
			if err := RDB.Del(ctx, key).Err(); err != nil {
				return err
			}
		}
		if cursor == 0 {
			break
		}
	}
	log.Logger.Debug("DelTeamsCache")
	return nil
}
