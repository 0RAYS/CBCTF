package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"github.com/vmihailenco/msgpack/v5"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

const rankingTTL = 30 * time.Second

// Rankings are complete, already ordered database snapshots. Atomic SET avoids
// both partial pipelines and floating-point composite-score precision loss.
func writeRanking[T any](ctx context.Context, key string, items []T) model.RetVal {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	data, err := msgpack.Marshal(items)
	if err == nil {
		err = RDB.Set(ctx, key, data, rankingTTL).Err()
	}
	if err != nil {
		return model.RetVal{Msg: i18n.Redis.SetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

func readRanking[T any](ctx context.Context, key string, start, end int64) ([]T, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	data, err := RDB.Get(ctx, key).Bytes()
	if errors.Is(err, goredis.Nil) {
		return nil, model.RetVal{Msg: i18n.Redis.NotFound}
	}
	var items []T
	if err == nil {
		err = msgpack.Unmarshal(data, &items)
	}
	if err != nil {
		return nil, model.RetVal{Msg: i18n.Redis.GetError, Attr: map[string]any{"Key": key, "Error": err.Error()}}
	}
	start = max(0, start)
	if end < start || start >= int64(len(items)) {
		return []T{}, model.SuccessRetVal()
	}
	return items[start:min(end+1, int64(len(items)))], model.SuccessRetVal()
}

func UpdateTeamRanking(ctx context.Context, contestID uint, teams []model.Team) model.RetVal {
	return writeRanking(ctx, fmt.Sprintf("ranking:contest:%d", contestID), teams)
}
func GetTeamRanking(ctx context.Context, contestID uint, start, end int64) ([]model.Team, model.RetVal) {
	return readRanking[model.Team](ctx, fmt.Sprintf("ranking:contest:%d", contestID), start, end)
}
func UpdateUserRanking(ctx context.Context, users []model.User) model.RetVal {
	return writeRanking(ctx, "ranking:users", users)
}
func GetUserRanking(ctx context.Context, start, end int64) ([]model.User, model.RetVal) {
	return readRanking[model.User](ctx, "ranking:users", start, end)
}
