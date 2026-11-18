package db

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/config"
	"CBCTF/internal/model"
)

// legacyContestChallengeTables 是迁移前的比赛题目表结构。
// 迁移需要在这套结构上跑，因此测试自建这套表，而不是直接用模型建表。
const legacyContestChallengeTables = `
CREATE TABLE challenges (
	id bigserial PRIMARY KEY,
	rand_id varchar(36) NOT NULL,
	name text,
	description text,
	category text,
	type text,
	generator_image text,
	network_policies jsonb,
	template jsonb DEFAULT '{}',
	created_at timestamptz,
	updated_at timestamptz,
	deleted_at timestamptz,
	version bigint
);
CREATE TABLE challenge_flags (
	id bigserial PRIMARY KEY,
	challenge_id bigint,
	value text,
	binding jsonb DEFAULT '{}',
	created_at timestamptz,
	updated_at timestamptz,
	deleted_at timestamptz,
	version bigint
);
CREATE TABLE contest_challenges (
	id bigserial PRIMARY KEY,
	contest_id bigint,
	challenge_id bigint,
	name text,
	description text,
	type text,
	category text,
	hidden boolean,
	attempt bigint,
	hints jsonb DEFAULT '[]',
	tags jsonb DEFAULT '[]',
	created_at timestamptz,
	updated_at timestamptz,
	deleted_at timestamptz,
	version bigint
);
CREATE TABLE contest_flags (
	id bigserial PRIMARY KEY,
	contest_id bigint,
	contest_challenge_id bigint,
	challenge_flag_id bigint,
	value text,
	score numeric DEFAULT 1000,
	current_score numeric DEFAULT 1000,
	decay numeric DEFAULT 50,
	min_score numeric DEFAULT 100,
	score_type bigint DEFAULT 0,
	solvers bigint,
	last timestamptz,
	created_at timestamptz,
	updated_at timestamptz,
	deleted_at timestamptz,
	version bigint
);
CREATE TABLE files (
	id bigserial PRIMARY KEY,
	model text NOT NULL,
	model_id bigint NOT NULL,
	rand_id varchar(36) NOT NULL,
	filename text,
	size bigint,
	path text,
	suffix text,
	hash text,
	type text,
	created_at timestamptz,
	updated_at timestamptz,
	deleted_at timestamptz,
	version bigint
);
CREATE TABLE generators (
	id bigserial PRIMARY KEY,
	challenge_id bigint,
	deleted_at timestamptz,
	created_at timestamptz,
	updated_at timestamptz,
	version bigint
);
CREATE TABLE victims (
	id bigserial PRIMARY KEY,
	challenge_id bigint,
	deleted_at timestamptz,
	created_at timestamptz,
	updated_at timestamptz,
	version bigint
);
`

// migrationTestDB 在独立 schema 中准备迁移前的表结构。
// 优先使用 CBCTF_TEST_POSTGRES_DSN，未设置时回退到项目自身的 config.yaml。
func migrationTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("CBCTF_TEST_POSTGRES_DSN")
	if dsn == "" {
		if _, err := os.Stat(filepath.Join("..", "..", "config.yaml")); err != nil {
			t.Skip("no test database configured")
		}
		config.Init(filepath.Join("..", "..", "config.yaml"))
		postgres := config.Env.Gorm.Postgres
		dsn = fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			postgres.Host, postgres.Port, postgres.User, postgres.Pwd, postgres.DB,
		)
		if config.Env != nil {
			// 附件复制只写入临时目录，不污染真实数据目录
			config.Env.Path = t.TempDir()
		}
	}
	admin, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Skipf("no test database available: %v", err)
	}
	pool, _ := admin.DB()
	t.Cleanup(func() { _ = pool.Close() })
	schema := fmt.Sprintf("cbctf_migration_test_%d", time.Now().UnixNano())
	if res := admin.Exec("CREATE SCHEMA " + schema); res.Error != nil {
		t.Skipf("cannot create test schema: %v", res.Error)
	}
	t.Cleanup(func() {
		cleanup, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
		if err != nil {
			return
		}
		defer func() {
			if pool, err := cleanup.DB(); err == nil {
				_ = pool.Close()
			}
		}()
		_ = cleanup.Exec("DROP SCHEMA IF EXISTS " + schema + " CASCADE").Error
	})
	database, err := gorm.Open(postgres.Open(dsn+" search_path="+schema), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	testPool, _ := database.DB()
	t.Cleanup(func() { _ = testPool.Close() })
	if res := database.Exec(legacyContestChallengeTables); res.Error != nil {
		t.Fatal(res.Error)
	}
	return database
}

func TestMigrateContestChallengeCopiesBackfillsLegacyRows(t *testing.T) {
	database := migrationTestDB(t)
	if res := database.Exec(
		`INSERT INTO challenges (rand_id, name, type, generator_image, network_policies, template) VALUES ('bank-rand', 'Bank', 'pods', 'registry/image:1', ?::jsonb, '{"pods":[{"key":"web","name":"web"}]}'::jsonb)`,
		`[{"service":"web"}]`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}
	if res := database.Exec(
		`INSERT INTO challenge_flags (challenge_id, value, binding) VALUES (1, 'flag{static:aaa}', '{"pod_key":"web","target":"FLAG"}'::jsonb)`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}
	if res := database.Exec(
		`INSERT INTO contest_challenges (contest_id, challenge_id, name, type) VALUES (7, 1, 'Bank', 'pods')`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}
	if res := database.Exec(
		`INSERT INTO contest_flags (contest_id, contest_challenge_id, challenge_flag_id, value) VALUES (7, 1, 1, 'flag{static:aaa}')`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}

	if ret := MigrateContestChallengeCopies(database); !ret.OK {
		t.Fatalf("migration failed: %+v", ret)
	}

	var contestChallenge model.ContestChallenge
	if res := database.First(&contestChallenge, 1); res.Error != nil {
		t.Fatal(res.Error)
	}
	if contestChallenge.RandID == "" || contestChallenge.RandID == "bank-rand" {
		t.Fatalf("contest challenge did not get its own rand id: %+v", contestChallenge)
	}
	if contestChallenge.GeneratorImage != "registry/image:1" {
		t.Fatalf("generator image was not copied: %+v", contestChallenge)
	}
	if len(contestChallenge.NetworkPolicies) != 1 || contestChallenge.NetworkPolicies[0].Service != "web" {
		t.Fatalf("network policies were not copied: %+v", contestChallenge.NetworkPolicies)
	}
	if len(contestChallenge.Template.Pods) != 1 || contestChallenge.Template.Pods[0].Key != "web" {
		t.Fatalf("template was not copied: %+v", contestChallenge.Template)
	}

	var contestFlag model.ContestFlag
	if res := database.First(&contestFlag, 1); res.Error != nil {
		t.Fatal(res.Error)
	}
	if contestFlag.Binding.PodKey != "web" || contestFlag.Binding.Target != "FLAG" {
		t.Fatalf("flag binding was not copied: %+v", contestFlag.Binding)
	}
	if !contestFlag.Active || contestFlag.Generation != 1 {
		t.Fatalf("flag generation was not initialized: %+v", contestFlag)
	}

	// 迁移必须幂等：再次执行不会覆盖已经迁移好的副本
	randID := contestChallenge.RandID
	if ret := MigrateContestChallengeCopies(database); !ret.OK {
		t.Fatalf("second migration failed: %+v", ret)
	}
	var again model.ContestChallenge
	if res := database.First(&again, 1); res.Error != nil {
		t.Fatal(res.Error)
	}
	if again.RandID != randID {
		t.Fatalf("rand id changed on rerun: %s -> %s", randID, again.RandID)
	}

	// 迁移后的数据必须满足 rand_id 唯一索引的要求
	if res := database.Exec(
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_contest_challenges_rand_id_active ON contest_challenges (rand_id) WHERE deleted_at IS NULL`,
	); res.Error != nil {
		t.Fatalf("rand id unique index failed on migrated rows: %v", res.Error)
	}
	if res := database.Exec(
		`INSERT INTO contest_challenges (contest_id, challenge_id, name, type, rand_id) VALUES (7, 1, 'Duplicate', 'pods', ?)`,
		randID,
	); res.Error == nil {
		t.Fatal("duplicated rand id was accepted")
	}
}

func TestAutoMigrateBuildsContestChallengeCopySchemaOnLegacyTable(t *testing.T) {
	database := migrationTestDB(t)
	if err := database.AutoMigrate(&model.ContestChallenge{}, &model.ContestFlag{}); err != nil {
		t.Fatalf("automigrate on legacy table failed: %v", err)
	}
	var indexCount int64
	if res := database.Raw(
		`SELECT COUNT(*) FROM pg_indexes WHERE tablename = 'contest_challenges' AND indexname = 'idx_contest_challenges_rand_id_active'`,
	).Scan(&indexCount); res.Error != nil {
		t.Fatal(res.Error)
	}
	if indexCount != 1 {
		t.Fatal("contest challenge rand id index was not created")
	}
}

func TestDeleteBankChallengeKeepsContestChallengeCopies(t *testing.T) {
	database := migrationTestDB(t)
	if res := database.Exec(
		`INSERT INTO challenges (rand_id, name, type) VALUES ('bank-rand', 'Bank', 'static')`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}
	if res := database.Exec(
		`INSERT INTO challenge_flags (challenge_id, value) VALUES (1, 'flag{static:aaa}')`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}
	if res := database.Exec(
		`INSERT INTO contest_challenges (contest_id, challenge_id, name, type) VALUES (7, 1, 'Bank', 'static')`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}
	if res := database.Exec(
		`INSERT INTO contest_flags (contest_id, contest_challenge_id, challenge_flag_id, value) VALUES (7, 1, 1, 'flag{static:aaa}')`,
	); res.Error != nil {
		t.Fatal(res.Error)
	}

	if ret := InitChallengeRepo(database).Delete("bank-rand"); !ret.OK {
		t.Fatalf("delete failed: %+v", ret)
	}

	var challengeCount, contestChallengeCount, contestFlagCount, challengeFlagCount int64
	database.Table("challenges").Where("deleted_at IS NULL").Count(&challengeCount)
	database.Table("contest_challenges").Where("deleted_at IS NULL").Count(&contestChallengeCount)
	database.Table("contest_flags").Where("deleted_at IS NULL").Count(&contestFlagCount)
	database.Table("challenge_flags").Where("deleted_at IS NULL").Count(&challengeFlagCount)
	if challengeCount != 0 || challengeFlagCount != 0 {
		t.Fatalf("bank challenge was not deleted: challenge=%d flag=%d", challengeCount, challengeFlagCount)
	}
	if contestChallengeCount != 1 || contestFlagCount != 1 {
		t.Fatalf("contest copies were deleted with the bank challenge: challenge=%d flag=%d", contestChallengeCount, contestFlagCount)
	}
}
