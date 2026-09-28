package db

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/model"
)

type prepullQuery struct {
	table string
	sql   string
	vars  []any
}

// Exercise the real repository's time filtering and generated PostgreSQL
// queries without needing an external database. Only contest rows are supplied.
func prepullQueryDB(t *testing.T, contests []model.Contest) (*gorm.DB, *[]prepullQuery) {
	t.Helper()
	tx, err := gorm.Open(postgres.Open("host=127.0.0.1 user=test dbname=test"), &gorm.Config{
		DryRun:               true,
		DisableAutomaticPing: true,
		Logger:               logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	pool, err := tx.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = pool.Close() })
	var queries []prepullQuery
	err = tx.Callback().Query().After("gorm:query").Register("test:prepull", func(query *gorm.DB) {
		// GORM also invokes query callbacks while building the nested subquery.
		// Only the two top-level FindAll calls have slice destinations.
		switch query.Statement.Dest.(type) {
		case *[]model.Contest, *[]model.Challenge:
		default:
			return
		}
		queries = append(queries, prepullQuery{
			table: query.Statement.Table,
			sql:   query.Statement.SQL.String(),
			vars:  append([]any(nil), query.Statement.Vars...),
		})
		if rows, ok := query.Statement.Dest.(*[]model.Contest); ok {
			*rows = append([]model.Contest(nil), contests...)
			query.RowsAffected = int64(len(contests))
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	return tx, &queries
}

func TestPrepullQueriesOnlyUnfinishedContestReferences(t *testing.T) {
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	contests := []model.Contest{
		{ID: 1, Start: now.Add(time.Hour), Duration: time.Hour},
		{ID: 2, Start: now.Add(-time.Hour), Duration: 2 * time.Hour},
		{ID: 3, Start: now.Add(-3 * time.Hour), Duration: time.Hour},
		{ID: 4, Start: now.Add(time.Hour), Duration: time.Hour, Hidden: true},
		{ID: 5, Start: now.Add(-time.Hour), Duration: time.Hour},
		{ID: 6, Start: now.Add(-time.Hour - time.Nanosecond), Duration: time.Hour},
	}
	tx, queries := prepullQueryDB(t, contests)
	if _, ret := InitChallengeRepo(tx).ListForUnfinishedContests(now); !ret.OK {
		t.Fatal(ret)
	}
	if len(*queries) != 2 {
		t.Fatalf("expected contest metadata and scoped challenge queries: %+v", *queries)
	}
	metadata, challenges := (*queries)[0], (*queries)[1]
	if metadata.table != "contests" || !strings.Contains(metadata.sql, `"contests"."deleted_at" IS NULL`) {
		t.Fatalf("contest query must exclude deleted contests: %s", metadata.sql)
	}
	if strings.Contains(metadata.sql, "SELECT *") {
		t.Fatalf("only contest IDs and schedule metadata are needed: %s", metadata.sql)
	}
	if !reflect.DeepEqual(challenges.vars, []any{uint(1), uint(2), uint(4), uint(5)}) {
		t.Fatalf("upcoming/running/end-time eligibility differs from Contest.IsOver: %v", challenges.vars)
	}
	for _, required := range []string{
		`id IN (SELECT "challenge_id" FROM "contest_challenges"`,
		`contest_id IN (`,
		`"contest_challenges"."deleted_at" IS NULL`,
		`"challenges"."deleted_at" IS NULL`,
	} {
		if !strings.Contains(challenges.sql, required) {
			t.Errorf("missing reference scope %q: %s", required, challenges.sql)
		}
	}
	if strings.Contains(challenges.sql, " JOIN ") || strings.Contains(challenges.sql, "hidden") {
		t.Fatalf("references should be deduplicated by membership and include hidden challenges: %s", challenges.sql)
	}
}

func TestPrepullSkipsChallengeQueryWithoutUnfinishedContests(t *testing.T) {
	now := time.Date(2026, 9, 28, 12, 0, 0, 0, time.UTC)
	for _, tc := range []struct {
		name     string
		contests []model.Contest
	}{
		{name: "no contests"},
		{
			name: "all ended",
			contests: []model.Contest{
				{ID: 1, Start: now.Add(-2 * time.Hour), Duration: time.Hour},
			},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			tx, queries := prepullQueryDB(t, tc.contests)
			challenges, ret := InitChallengeRepo(tx).ListForUnfinishedContests(now)
			if !ret.OK || len(challenges) != 0 {
				t.Fatalf("unexpected challenges: %v, result: %+v", challenges, ret)
			}
			if len(*queries) != 1 || (*queries)[0].table != "contests" {
				t.Fatalf("empty eligibility must not scan challenges: %+v", *queries)
			}
		})
	}
}
