package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func RecordDockerCreate(teamID uint, challengeID string) error {
	ctx := context.Background()
	err := RDB.Set(ctx, fmt.Sprintf("d:c:%d:%s", teamID, challengeID), "1", 1*time.Minute).Err()
	return err
}

func CheckDockerCreate(teamID uint, challengeID string) (bool, error) {
	ctx := context.Background()
	_, err := RDB.Get(ctx, fmt.Sprintf("d:c:%d:%s", teamID, challengeID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
