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
		return model.FailRetVal(i18n.Redis.SetError, "Cache", key, err)
	}
	return model.SuccessRetVal()
}

func readRanking[T any](ctx context.Context, key string) ([]T, model.RetVal) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	data, err := RDB.Get(ctx, key).Bytes()
	if errors.Is(err, goredis.Nil) {
		return nil, model.FailRetVal(i18n.Redis.NotFound, "Cache", key, nil)
	}
	var items []T
	if err == nil {
		err = msgpack.Unmarshal(data, &items)
	}
	if err != nil {
		return nil, model.FailRetVal(i18n.Redis.GetError, "Cache", key, err)
	}
	return items, model.SuccessRetVal()
}

func UpdateTeamRanking(ctx context.Context, contestID uint, teams []model.Team) model.RetVal {
	return writeRanking(ctx, fmt.Sprintf("ranking:contest:%d", contestID), teams)
}
// GetTeamRanking 返回完整的队伍排名快照，隐藏队伍需要在分页前剔除。
func GetTeamRanking(ctx context.Context, contestID uint) ([]model.Team, model.RetVal) {
	return readRanking[model.Team](ctx, fmt.Sprintf("ranking:contest:%d", contestID))
}
func UpdateUserRanking(ctx context.Context, users []model.User) model.RetVal {
	return writeRanking(ctx, "ranking:users", users)
}
// GetUserRanking 返回完整的用户排名快照，隐藏用户需要在分页前剔除。
func GetUserRanking(ctx context.Context) ([]model.User, model.RetVal) {
	return readRanking[model.User](ctx, "ranking:users")
}
