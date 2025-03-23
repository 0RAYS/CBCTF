package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// RecordContainersCreate 记录 1*time.Minutes 内启动容器
func RecordContainersCreate(teamID uint, challengeID string) error {
	ctx := context.Background()
	err := RDB.Set(ctx, fmt.Sprintf("d:c:%d:%s", teamID, challengeID), "1", 1*time.Minute).Err()
	return err
}

// CheckContainersCreate 是否 1*time.Minutes 内启动过容器
func CheckContainersCreate(teamID uint, challengeID string) (bool, error) {
	ctx := context.Background()
	_, err := RDB.Get(ctx, fmt.Sprintf("d:c:%d:%s", teamID, challengeID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
