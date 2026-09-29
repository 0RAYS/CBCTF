package service

import (
	"errors"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/db"
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func queryFailureDB(t *testing.T) *gorm.DB {
	t.Helper()
	log.Init()
	database, err := gorm.Open(postgres.Open("host=127.0.0.1 dbname=test"), &gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := database.DB()
	t.Cleanup(func() { _ = pool.Close() })
	return database
}

func TestQueryFailuresDoNotBecomeBusinessStates(t *testing.T) {
	database := queryFailureDB(t)
	if err := database.Callback().Query().Before("gorm:query").Register("test:query_failure", func(tx *gorm.DB) { _ = tx.AddError(errors.New("database unavailable")) }); err != nil {
		t.Fatal(err)
	}
	flags := []model.ContestFlag{{ID: 1, ContestChallengeID: 1}}
	for _, check := range []func() (bool, model.RetVal){
		func() (bool, model.RetVal) { return CheckIfGenerated(database, model.Team{ID: 1}, flags) },
		func() (bool, model.RetVal) { return CheckIfSolved(database, model.Team{ID: 1}, flags) },
		func() (bool, model.RetVal) {
			return db.InitContestChallengeRepo(database).IsUniqueContestChallenge(1, 1)
		},
	} {
		if value, ret := check(); ret.OK || value {
			t.Fatalf("query failure became successful state: %v %+v", value, ret)
		}
	}
	if _, ret := CountAttempts(database, model.Team{ID: 1}, model.ContestChallenge{ID: 1}); ret.OK {
		t.Fatal("query failure became zero attempts")
	}
	created := false
	if err := database.Callback().Create().Before("gorm:create").Register("test:no_create", func(*gorm.DB) { created = true }); err != nil {
		t.Fatal(err)
	}
	_, ret := CreateTeamFlag(database, model.Team{ID: 1}, model.Contest{}, model.ContestChallenge{ContestFlags: flags})
	if ret.OK || created {
		t.Fatalf("failed lookup triggered flag creation: %+v", ret)
	}
}

func TestFlagLookupFailureIsNotWrongSubmission(t *testing.T) {
	database := queryFailureDB(t)
	if err := database.Callback().Query().Before("gorm:query").Register("test:flag_lookup", func(tx *gorm.DB) {
		if tx.Statement.Schema.Name == "TeamFlag" {
			_ = tx.AddError(errors.New("read failed"))
			return
		}
		switch value := tx.Statement.Dest.(type) {
		case *int64:
			*value = 1
		case *[]model.ContestFlag:
			*value = []model.ContestFlag{{ID: 1}}
		}
		tx.RowsAffected = 1
	}); err != nil {
		t.Fatal(err)
	}
	_, _, _, ret := VerifyFlag(database, model.Team{ID: 1}, model.ContestChallenge{ID: 1}, "flag{value}")
	if ret.OK || ret.Msg != i18n.Model.TeamFlag.GetError {
		t.Fatalf("read failed but flag treated as mismatch: %+v", ret)
	}
}
