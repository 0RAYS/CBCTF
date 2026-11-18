package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"CBCTF/internal/model"
)

// imageDetectionTTL keeps the last detection result of a scope for a day.
const imageDetectionTTL = 24 * time.Hour

func imageDetectionKey(scope string) string {
	return fmt.Sprintf("images:detect:%s", scope)
}

// SaveImageDetection keeps only the last detection result of a scope.
func SaveImageDetection(ctx context.Context, scope string, detection model.ImageDetection) error {
	data, err := json.Marshal(detection)
	if err != nil {
		return err
	}
	return RDB.Set(ctx, imageDetectionKey(scope), data, imageDetectionTTL).Err()
}

// LoadImageDetection returns the last detection result of a scope, if any.
func LoadImageDetection(ctx context.Context, scope string) (*model.ImageDetection, error) {
	data, err := RDB.Get(ctx, imageDetectionKey(scope)).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var detection model.ImageDetection
	if err := json.Unmarshal(data, &detection); err != nil {
		return nil, err
	}
	return &detection, nil
}
