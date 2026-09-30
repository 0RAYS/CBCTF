package db

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"CBCTF/internal/traffic"
)

// Temporary tables shadow production table names only inside this transaction.
func TestTrafficEvidencePostgres(t *testing.T) {
	dsn := os.Getenv("CBCTF_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("CBCTF_TEST_POSTGRES_DSN is not set")
	}
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	pool, err := database.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer func(pool *sql.DB) {
		_ = pool.Close()
	}(pool)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	tx := database.WithContext(ctx).Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	defer tx.Rollback()
	for _, ddl := range []string{
		`CREATE TEMP TABLE traffics (id bigserial PRIMARY KEY, victim_id bigint UNIQUE, ips jsonb, accesses jsonb, analysis jsonb, archived boolean DEFAULT false, created_at timestamptz, updated_at timestamptz, deleted_at timestamptz, version bigint) ON COMMIT DROP`,
		`CREATE TEMP TABLE victims (id bigint PRIMARY KEY, team_id bigint, contest_id bigint, created_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE teams (id bigint PRIMARY KEY, deleted_at timestamptz) ON COMMIT DROP`,
		`INSERT INTO teams(id) VALUES (1),(2),(3)`,
		`INSERT INTO victims(id,team_id,contest_id,created_at) VALUES (1,1,10,'2000-01-01'),(2,2,10,'2000-01-01'),(3,3,11,'2000-01-01')`,
	} {
		if err = tx.Exec(ddl).Error; err != nil {
			t.Fatal(err)
		}
	}
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	repo := InitTrafficRepo(tx)
	for _, id := range []uint{1, 2, 3} {
		accesses := []traffic.Access{{IP: "8.8.8.8", Time: at, Source: "tcp_syn"}, {IP: "1.1.1.1", Time: at.Add(-time.Hour), Source: "tcp_syn"}}
		if ret := repo.ReplaceAnalysis(id, []string{"10.0.0.2"}, accesses, &traffic.AnalysisReport{}, false); !ret.OK {
			t.Fatalf("save: %+v", ret)
		}
	}
	rows, ret := repo.ListSharedContestVictimIPs(10, at, at.Add(time.Hour))
	if !ret.OK || len(rows) != 2 {
		t.Fatalf("contest/time-scoped overlap: %+v %+v", rows, ret)
	}
	for _, row := range rows {
		if row.SrcIP != "8.8.8.8" || !row.FirstTime.Equal(at) || len(row.VictimIDs) != 1 {
			t.Fatalf("incorrect evidence: %+v", row)
		}
	}
	if ret = repo.ReplaceAnalysis(1, []string{}, []traffic.Access{}, &traffic.AnalysisReport{}, true); !ret.OK {
		t.Fatalf("archive: %+v", ret)
	}
	if ret = repo.ReplaceAnalysis(1, []string{"8.8.8.8"}, []traffic.Access{}, nil, false); !ret.OK {
		t.Fatalf("live write: %+v", ret)
	}
	record, ret := repo.GetAnalysis(1)
	if !ret.OK || !record.Archived || record.Analysis == nil || len(record.IPs) != 0 {
		t.Fatalf("live snapshot overwrote archive: %+v %+v", record, ret)
	}
}
