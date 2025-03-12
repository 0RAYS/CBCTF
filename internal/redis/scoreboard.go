package redis

import (
	"CBCTF/internal/config"
	"CBCTF/internal/model"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"
	"time"
)

func UpdateTeamRanking(contestID uint, teams []model.Team) error {
	if !config.Env.Redis.On {
		return nil
	}
	key := fmt.Sprintf("%d:rank", contestID)
	ctx := context.Background()
	pipe := RDB.Pipeline()
	pipe.Del(ctx, key)

	for _, team := range teams {
		timestamp := team.Last.UnixMilli()
		compositeScore := team.Score*1e13 + float64(1e13-timestamp)
		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  compositeScore,
			Member: team.ID,
		})
		data, _ := msgpack.Marshal(&team)
		pipe.Set(ctx, fmt.Sprintf("team:%d", team.ID), data, 5*time.Minute)
	}
	pipe.Expire(ctx, key, 5*time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}

func GetTeamRanking(contestID uint, start int64, end int64) ([]model.Team, error) {
	if !config.Env.Redis.On {
		return make([]model.Team, 0), nil
	}
	key := fmt.Sprintf("%d:rank", contestID)
	ctx := context.Background()
	results, err := RDB.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return make([]model.Team, 0), err
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
		err := msgpack.Unmarshal(str, &t)
		if err != nil {
			return make([]model.Team, 0), err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func UpdateUserRanking(users []model.User) error {
	if !config.Env.Redis.On {
		return nil
	}
	key := "users:rank"
	ctx := context.Background()
	pipe := RDB.Pipeline()
	pipe.Del(ctx, key)
	for _, user := range users {
		compositeScore := user.Score*1e5 + float64(user.Solved)
		pipe.ZAdd(ctx, key, &redis.Z{
			Score:  compositeScore,
			Member: user.ID,
		})
		data, _ := msgpack.Marshal(&user)
		pipe.Set(ctx, fmt.Sprintf("user:%d", user.ID), data, 12*time.Hour)
	}
	pipe.Expire(ctx, key, 12*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func GetUserRanking(start int64, end int64) ([]model.User, error) {
	if !config.Env.Redis.On {
		return make([]model.User, 0), nil
	}
	key := "users:rank"
	ctx := context.Background()
	results, err := RDB.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return make([]model.User, 0), err
	}
	pipe := RDB.Pipeline()
	for _, res := range results {
		userID := res.Member.(string)
		pipe.Get(ctx, fmt.Sprintf("user:%s", userID))
	}
	cmds, _ := pipe.Exec(ctx)
	var users []model.User
	for _, cmd := range cmds {
		str, _ := cmd.(*redis.StringCmd).Bytes()
		var u model.User
		err := msgpack.Unmarshal(str, &u)
		if err != nil {
			return make([]model.User, 0), err
		}
		users = append(users, u)
	}
	return users, nil
}
