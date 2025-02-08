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
	TeamPattern  = "t:%d:%v" // team:<team_id>:<preload>
	TeamsPattern = "ts:%v"   // teams:<preload>
)

func GetTeamCache(id uint, preload int) (model.Team, bool) {
	if !config.Env.Redis.On {
		return model.Team{}, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf(TeamPattern, id, preload)).Result()
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

func GetTeamsCache(preload int) ([]model.Team, bool) {
	if !config.Env.Redis.On {
		return nil, false
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf(TeamsPattern, preload)).Result()
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

func SetTeamCache(team model.Team, preload int) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(team)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, fmt.Sprintf(TeamPattern, team.ID, preload), data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetTeamCache: ", team.ID)
	return nil
}

func SetTeamsCache(teams []model.Team, preload int) error {
	if !config.Env.Redis.On {
		return nil
	}
	ctx := context.Background()
	data, err := msgpack.Marshal(teams)
	if err != nil {
		return err
	}
	if err = RDB.Set(ctx, fmt.Sprintf(TeamsPattern, preload), data, time.Minute).Err(); err != nil {
		return err
	}
	log.Logger.Debug("SetTeamsCache: ", len(teams))
	return nil
}

func DelTeamCache(id uint) error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(fmt.Sprintf(TeamPattern, id, "*"))
}

func DelTeamsCache() error {
	if !config.Env.Redis.On {
		return nil
	}
	return DeleteKeysByPattern(fmt.Sprintf(TeamsPattern, "*"))
}
