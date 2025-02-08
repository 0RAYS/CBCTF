package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func UpdateRanking(contestID uint, teams []model.Team) error {
	if !config.Env.Redis.On {
		return nil
	}
	key := fmt.Sprintf("%d:rank", contestID)
	ctx := context.Background()
	pipe := RDB.Pipeline()
	pipe.Del(ctx, key)

	for _, team := range teams {
		timestamp := team.Last.UnixNano()
		compositeScore := float64(team.Score)*1e13 + float64(1e18-timestamp)
		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  compositeScore,
			Member: team.ID,
		})
		data, _ := json.Marshal(team)
		pipe.Set(ctx, fmt.Sprintf("team:%d", team.ID), data, 1*time.Hour)
	}

	pipe.Expire(ctx, key, time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}

func GetCachedRanking(contestID uint, limit int64, offset int64) ([]model.Team, error) {
	if !config.Env.Redis.On {
		return nil, nil
	}
	key := fmt.Sprintf("%d:rank", contestID)
	ctx := context.Background()
	results, err := RDB.ZRevRangeWithScores(ctx, key, offset, limit).Result()
	if err != nil {
		return nil, err
	}

	pipe := RDB.Pipeline()
	for _, res := range results {
		teamID := res.Member.(string)
		pipe.Get(ctx, fmt.Sprintf("team:%s", teamID))
	}
	cmds, _ := pipe.Exec(ctx)

	var teams []model.Team
	for _, cmd := range cmds {
		str, _ := cmd.(*redis.StringCmd).Bytes()
		var t model.Team
		err := json.Unmarshal(str, &t)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}
