package service

import (
	"strings"
	"testing"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"CBCTF/internal/model"
	"CBCTF/internal/redis"
)

func TestHiddenTeamsAreFilteredBeforePaginationAndRanksStayContiguous(t *testing.T) {
	database := queryFailureDB(t)
	old := redis.RDB
	client := goredis.NewClient(&goredis.Options{})
	redis.RDB = client
	t.Cleanup(func() { redis.RDB = old; _ = client.Close() })
	client.AddHook(&unavailableRankingCache{})

	var statements []string
	_ = database.Callback().Query().After("gorm:query").Register("test:ranking_rows", func(tx *gorm.DB) {
		statements = append(statements, tx.Statement.SQL.String())
		switch dst := tx.Statement.Dest.(type) {
		case *int64:
			*dst = 2
		case *[]model.Team:
			*dst = []model.Team{{ID: 1, Score: 300}, {ID: 2, Score: 200}}
		}
		tx.RowsAffected = 2
	})

	teams, count, ret := GetTeamRanking(database, model.Contest{ID: 1}, 10, 0, false)
	if !ret.OK || count != 2 || len(teams) != 2 {
		t.Fatalf("unexpected ranking result: teams=%+v count=%d ret=%+v", teams, count, ret)
	}
	if teams[0].Rank != 1 || teams[1].Rank != 2 {
		t.Fatalf("ranks are not contiguous: %+v", teams)
	}
	if !strings.Contains(strings.Join(statements, "\n"), "hidden") {
		t.Fatalf("hidden teams are not filtered in SQL: %v", statements)
	}

	statements = nil
	teams, count, ret = GetTeamRanking(database, model.Contest{ID: 1}, 10, 0, true)
	if !ret.OK || count != 2 || len(teams) != 2 {
		t.Fatalf("unexpected ranking result: teams=%+v count=%d ret=%+v", teams, count, ret)
	}
	if strings.Contains(strings.Join(statements, "\n"), "hidden") {
		t.Fatalf("hidden teams must be listed when requested: %v", statements)
	}
}
