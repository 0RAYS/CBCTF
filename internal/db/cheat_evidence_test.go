package db

import (
	"context"
	"os"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"CBCTF/internal/traffic"
)

// cheatEvidenceDB 在事务内建立临时表，模拟一次比赛的可观测数据。
func cheatEvidenceDB(t *testing.T) *gorm.DB {
	t.Helper()
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
	t.Cleanup(func() { _ = pool.Close() })
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	t.Cleanup(cancel)
	tx := database.WithContext(ctx).Begin()
	if tx.Error != nil {
		t.Fatal(tx.Error)
	}
	t.Cleanup(func() { _ = tx.Rollback() })
	for _, ddl := range []string{
		`CREATE TEMP TABLE users (id bigint PRIMARY KEY, deleted_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE user_contests (user_id bigint, contest_id bigint) ON COMMIT DROP`,
		`CREATE TEMP TABLE teams (id bigint PRIMARY KEY, contest_id bigint, deleted_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE requests (id bigserial PRIMARY KEY, ip text, user_id bigint, time timestamptz, method text, path text, status int, latency bigint, user_agent text, deleted_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE submissions (id bigserial PRIMARY KEY, team_id bigint, user_id bigint, contest_id bigint, contest_challenge_id bigint, contest_flag_id bigint, ip text, created_at timestamptz, value text, solved boolean, deleted_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE contest_challenges (id bigint PRIMARY KEY, name text, deleted_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE contest_flags (id bigint PRIMARY KEY, contest_challenge_id bigint, deleted_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE team_flags (id bigint PRIMARY KEY, team_id bigint, contest_flag_id bigint, value text, deleted_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE victims (id bigint PRIMARY KEY, team_id bigint, contest_id bigint, contest_challenge_id bigint, deleted_at timestamptz, created_at timestamptz) ON COMMIT DROP`,
		`CREATE TEMP TABLE traffics (id bigserial PRIMARY KEY, victim_id bigint UNIQUE, ips jsonb, accesses jsonb, analysis jsonb, archived boolean DEFAULT false, created_at timestamptz, updated_at timestamptz, deleted_at timestamptz, version bigint) ON COMMIT DROP`,
		`INSERT INTO users(id) VALUES (1),(2),(3)`,
		`INSERT INTO user_contests(user_id, contest_id) VALUES (1,10),(2,10),(3,10)`,
		`INSERT INTO teams(id, contest_id) VALUES (1,10),(2,10)`,
		`INSERT INTO contest_challenges(id, name) VALUES (100,'web'),(200,'pwn')`,
	} {
		if err = tx.Exec(ddl).Error; err != nil {
			t.Fatal(err)
		}
	}
	return tx
}

func TestSharedContestUserIPsTraceRequests(t *testing.T) {
	tx := cheatEvidenceDB(t)
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	for _, row := range []struct {
		userID uint
		offset time.Duration
		path   string
	}{
		{1, 0, "/api/challenges"},
		{2, 3 * time.Second, "/api/submit"},
		{1, 10 * time.Second, "/api/challenges"},
		{2, 12 * time.Second, "/api/submit"},
	} {
		statement := `INSERT INTO requests(ip, user_id, time, method, path, status, latency, user_agent)
			VALUES ('8.8.8.8', ?, ?, 'POST', ?, 200, 12, 'agent')`
		if err := tx.Exec(statement, row.userID, at.Add(row.offset), row.path).Error; err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Exec(`INSERT INTO requests(ip, user_id, time, method, path, status) VALUES ('1.1.1.1', 3, ?, 'GET', '/', 200)`, at).Error; err != nil {
		t.Fatal(err)
	}

	repo := InitRequestRepo(tx)
	rows, ret := repo.ListSharedContestUserIPs(10, at.Add(-time.Hour), at.Add(time.Hour))
	if !ret.OK {
		t.Fatalf("shared IPs: %+v", ret)
	}
	if len(rows) != 2 {
		t.Fatalf("expected two users on the shared IP, got %+v", rows)
	}
	if rows[0].IP != "8.8.8.8" || rows[0].TotalRequests != 4 || rows[0].TotalUsers != 2 {
		t.Fatalf("wrong shared IP aggregate: %+v", rows[0])
	}
	if rows[0].MinGapSeconds != 3 {
		t.Fatalf("expected a 3s gap between different users, got %+v", rows[0])
	}
	if !rows[0].FirstSeen.Equal(at) || !rows[0].LastSeen.Equal(at.Add(12*time.Second)) {
		t.Fatalf("wrong observation window: %+v", rows[0])
	}

	requests, ret := repo.ListContestIPRequests(10, []string{"8.8.8.8"}, at.Add(-time.Hour), at.Add(time.Hour), 3)
	if !ret.OK || len(requests) != 3 {
		t.Fatalf("request detail: %+v %+v", requests, ret)
	}
	if requests[0].Path != "/api/challenges" || requests[0].UserAgent != "agent" {
		t.Fatalf("wrong request detail: %+v", requests[0])
	}
}

func TestSharedContestUserIPsTraceSubmissions(t *testing.T) {
	tx := cheatEvidenceDB(t)
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	statement := `INSERT INTO submissions(team_id, user_id, contest_id, contest_challenge_id, ip, created_at, value, solved)
		VALUES (1, 1, 10, 100, '8.8.8.8', ?, 'flag{guess}', false)`
	if err := tx.Exec(statement, at).Error; err != nil {
		t.Fatal(err)
	}

	submissions, ret := InitSubmissionRepo(tx).ListContestIPSubmissions(10, []string{"8.8.8.8"}, at.Add(-time.Hour), at.Add(time.Hour), 10)
	if !ret.OK || len(submissions) != 1 {
		t.Fatalf("submission detail: %+v %+v", submissions, ret)
	}
	if submissions[0].ChallengeName != "web" || submissions[0].TeamID != 1 || submissions[0].Value != "flag{guess}" {
		t.Fatalf("wrong submission detail: %+v", submissions[0])
	}
}

func TestSharedContestVictimIPsTraceAccesses(t *testing.T) {
	tx := cheatEvidenceDB(t)
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	for _, victim := range []struct {
		id        uint
		teamID    uint
		challenge uint
		offset    time.Duration
	}{
		{1, 1, 100, 0},
		{2, 2, 200, 5 * time.Second},
	} {
		statement := `INSERT INTO victims(id, team_id, contest_id, contest_challenge_id, created_at) VALUES (?, ?, 10, ?, ?)`
		if err := tx.Exec(statement, victim.id, victim.teamID, victim.challenge, at).Error; err != nil {
			t.Fatal(err)
		}
		accesses := []traffic.Access{{IP: "8.8.8.8", Time: at.Add(victim.offset), Source: "tcp_syn"}}
		if ret := InitTrafficRepo(tx).ReplaceAnalysis(victim.id, []string{"10.0.0.2"}, accesses, &traffic.AnalysisReport{}, false); !ret.OK {
			t.Fatalf("save traffic: %+v", ret)
		}
	}

	repo := InitTrafficRepo(tx)
	rows, ret := repo.ListSharedContestVictimIPs(10, at.Add(-time.Hour), at.Add(time.Hour))
	if !ret.OK || len(rows) != 2 {
		t.Fatalf("shared victim IPs: %+v %+v", rows, ret)
	}
	if rows[0].TeamID != 1 || rows[0].AccessCount != 1 || rows[0].MinGapSeconds != 5 {
		t.Fatalf("wrong victim aggregate: %+v", rows[0])
	}

	accesses, ret := repo.ListContestVictimAccesses(10, []string{"8.8.8.8"}, at.Add(-time.Hour), at.Add(time.Hour), 10)
	if !ret.OK || len(accesses) != 2 {
		t.Fatalf("victim access detail: %+v %+v", accesses, ret)
	}
	if accesses[1].TeamID != 2 || accesses[1].ChallengeName != "pwn" {
		t.Fatalf("wrong victim access detail: %+v", accesses[1])
	}
}

func TestWrongFlagSubmissionsTraceOwner(t *testing.T) {
	tx := cheatEvidenceDB(t)
	at := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	for _, ddl := range []string{
		`INSERT INTO contest_flags(id, contest_challenge_id) VALUES (300,100)`,
		`INSERT INTO team_flags(id, team_id, contest_flag_id, value) VALUES (400,2,300,'flag{real}')`,
	} {
		if err := tx.Exec(ddl).Error; err != nil {
			t.Fatal(err)
		}
	}
	for _, statement := range []string{
		`INSERT INTO submissions(team_id, user_id, contest_id, contest_challenge_id, contest_flag_id, ip, created_at, value, solved)
			VALUES (2,2,10,100,300,'9.9.9.9',?,'flag{real}',true)`,
		`INSERT INTO submissions(team_id, user_id, contest_id, contest_challenge_id, ip, created_at, value, solved)
			VALUES (1,1,10,100,'8.8.8.8',?,'flag{real}',false)`,
	} {
		if err := tx.Exec(statement, at).Error; err != nil {
			t.Fatal(err)
		}
	}

	rows, ret := InitTeamFlagRepo(tx).ListContestWrongFlagSubmissions(10, at.Add(-time.Hour), at.Add(time.Hour))
	if !ret.OK || len(rows) != 1 {
		t.Fatalf("wrong flag submissions: %+v %+v", rows, ret)
	}
	if rows[0].MatchedTeamID != 2 || rows[0].TeamID != 1 || rows[0].ChallengeName != "web" {
		t.Fatalf("wrong wrong-flag row: %+v", rows[0])
	}
	if !rows[0].MatchedSolvedAt.Equal(at) {
		t.Fatalf("expected the legitimate solve time, got %+v", rows[0].MatchedSolvedAt)
	}
}
