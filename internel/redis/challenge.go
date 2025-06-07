package redis

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// RecordChallengeInit 记录 1*time.Minutes 内初始化
func RecordChallengeInit(teamID uint, contestChallengeID uint) error {
	ctx := context.Background()
	err := RDB.Set(ctx, fmt.Sprintf("c:init:%d:%d", teamID, contestChallengeID), "1", 1*time.Minute).Err()
	return err
}

// CheckChallengeInit 检查 1*time.Minutes 内是否初始化
func CheckChallengeInit(teamID uint, contestChallengeID uint) (bool, error) {
	ctx := context.Background()
	_, err := RDB.Get(ctx, fmt.Sprintf("c:init:%d:%d", teamID, contestChallengeID)).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
