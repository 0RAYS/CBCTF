package redis

import (
	"CBCTF/internal/model"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	teamRankingKey = "contests:%d:teams:rank"
	userRankingKey = "users:rank"
)

func UpdateTeamRanking(contestID uint, teams []model.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(teamRankingKey, contestID)
	pipe := RDB.Pipeline()
	pipe.Del(ctx, key)

	for _, team := range teams {
		timestamp := team.Last.UnixMilli()
		compositeScore := team.Score*1e16 + float64(1e13-timestamp)
		pipe.ZAdd(ctx, key, redis.Z{
			Score:  compositeScore,
			Member: fmt.Sprintf(teamKey, contestID, team.ID),
		})
		data, _ := msgpack.Marshal(&team)
		pipe.Set(ctx, fmt.Sprintf(teamKey, contestID, team.ID), data, 5*time.Minute)
	}
	pipe.Expire(ctx, key, 5*time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}

func GetTeamRanking(contestID uint, start int64, end int64) ([]model.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	key := fmt.Sprintf(teamRankingKey, contestID)
	teams := make([]model.Team, 0)
	results, err := RDB.ZRevRangeWithScores(ctx, key, start, end).Result()
	if err != nil {
		return nil, err
	}

	pipe := RDB.Pipeline()
	for _, res := range results {
		pipe.Get(ctx, res.Member.(string))
	}
	cmds, _ := pipe.Exec(ctx)

	for _, cmd := range cmds {
		str, _ := cmd.(*redis.StringCmd).Bytes()
		var t model.Team
		err = msgpack.Unmarshal(str, &t)
		if err != nil {
			return nil, err
		}
		teams = append(teams, t)
	}
	return teams, nil
}

func UpdateUserRanking(users []model.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	pipe := RDB.Pipeline()
	pipe.Del(ctx, userRankingKey)
	for _, user := range users {
		compositeScore := user.Score*1e8 + float64(user.Solved)
		pipe.ZAdd(ctx, userRankingKey, redis.Z{
			Score:  compositeScore,
			Member: fmt.Sprintf(userKey, user.ID),
		})
		data, _ := msgpack.Marshal(&user)
		pipe.Set(ctx, fmt.Sprintf(userKey, user.ID), data, 12*time.Hour)
	}
	pipe.Expire(ctx, userRankingKey, 12*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func GetUserRanking(start int64, end int64) ([]model.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	results, err := RDB.ZRevRangeWithScores(ctx, userRankingKey, start, end).Result()
	if err != nil {
		return nil, err
	}
	pipe := RDB.Pipeline()
	for _, res := range results {
		pipe.Get(ctx, res.Member.(string))
	}
	cmds, _ := pipe.Exec(ctx)
	var users []model.User
	for _, cmd := range cmds {
		str, _ := cmd.(*redis.StringCmd).Bytes()
		var u model.User
		err = msgpack.Unmarshal(str, &u)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}
