package redis

import (
	"CBCTF/internal/config"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func RecordChallengeInit(teamID uint, challengeID string) error {
	if !config.Env.Redis.On {
		return errors.New("RedisOff")
	}
	ctx := context.Background()
	err := RDB.Set(ctx, fmt.Sprintf("c:i:%d:%s", teamID, challengeID), "1", 1*time.Minute).Err()
	return err
}

func CheckChallengeInit(teamID uint, challengeID string) (bool, error) {
	if !config.Env.Redis.On {
		return false, errors.New("RedisOff")
	}
	ctx := context.Background()
	data, err := RDB.Get(ctx, fmt.Sprintf("c:i:%d:%s", teamID, challengeID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return data == "1", nil
}
