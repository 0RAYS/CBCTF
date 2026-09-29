package service

import (
	"context"
	"errors"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"CBCTF/internal/model"
	"CBCTF/internal/redis"
)

type unavailableRankingCache struct{ calls int }

func (*unavailableRankingCache) DialHook(next goredis.DialHook) goredis.DialHook { return next }
func (*unavailableRankingCache) ProcessPipelineHook(next goredis.ProcessPipelineHook) goredis.ProcessPipelineHook {
	return next
}
func (h *unavailableRankingCache) ProcessHook(_ goredis.ProcessHook) goredis.ProcessHook {
	return func(context.Context, goredis.Cmder) error { h.calls++; return errors.New("cache unavailable") }
}

func TestRankingCacheFailureReadsDatabaseWithoutRebuild(t *testing.T) {
	database := queryFailureDB(t)
	old := redis.RDB
	client := goredis.NewClient(&goredis.Options{})
	redis.RDB = client
	t.Cleanup(func() { redis.RDB = old; _ = client.Close() })
	hook := &unavailableRankingCache{}
	client.AddHook(hook)
	writes := 0
	_ = database.Callback().Update().Before("gorm:update").Register("test:no_update", func(*gorm.DB) { writes++ })
	_ = database.Callback().Query().After("gorm:query").Register("test:database_rank", func(tx *gorm.DB) {
		switch dst := tx.Statement.Dest.(type) {
		case *int64:
			*dst = 2
		case *[]model.Team:
			*dst = []model.Team{{ID: 2, Score: 100}, {ID: 1, Score: 90}}
		}
		tx.RowsAffected = 2
	})
	teams, count, ret := GetTeamRanking(database, model.Contest{ID: 1}, 10, 0)
	if !ret.OK || count != 2 || len(teams) != 2 || teams[0].ID != 2 || teams[0].Rank != 1 || writes != 0 || hook.calls != 1 {
		t.Fatalf("cache error caused rebuild or lost DB data: %+v count=%d ret=%+v writes=%d cache_calls=%d", teams, count, ret, writes, hook.calls)
	}
}
