package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// RecordVictimCreate 记录 1*time.Minutes 内启动靶机
func RecordVictimCreate(teamID uint, contestChallengeID uint) error {
	ctx := context.Background()
	err := RDB.Set(ctx, fmt.Sprintf("victim:%d:%d", teamID, contestChallengeID), "1", 1*time.Minute).Err()
	return err
}

// CheckVictimCreate 是否 1*time.Minutes 内启动过靶机
func CheckVictimCreate(teamID uint, contestChallengeID uint) (bool, error) {
	ctx := context.Background()
	_, err := RDB.Get(ctx, fmt.Sprintf("victim:%d:%d", teamID, contestChallengeID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
