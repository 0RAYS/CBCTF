package cron

import (
	"errors"
	"sync"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

func TestMutexCleanupKeepsLocksOnDatabaseFailure(t *testing.T) {
	old, oldLogger := db.CronDB, log.Logger
	t.Cleanup(func() { db.CronDB, log.Logger = old, oldLogger })
	log.Init()
	database, err := gorm.Open(postgres.Open("host=127.0.0.1 dbname=test"), &gorm.Config{DryRun: true, DisableAutomaticPing: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := database.DB()
	defer pool.Close()
	db.CronDB = database
	_ = database.Callback().Query().Before("gorm:query").Register("test:offline", func(tx *gorm.DB) { _ = tx.AddError(errors.New("database offline")) })
	for _, tc := range []struct {
		pool *sync.Map
		key  any
		run  func() model.RetVal
	}{
		{&service.SolvedMutex, uint(987654), clearSubmissionMutexTask},
		{&db.CheatMutex, "test-cleanup-key", clearCheatMutexTask},
		{&service.JoinTeamMutex, uint(987654), clearJoinTeamMutexTask},
	} {
		lock := &sync.Mutex{}
		tc.pool.Store(tc.key, lock)
		ret := tc.run()
		stored, exists := tc.pool.Load(tc.key)
		tc.pool.Delete(tc.key)
		if ret.OK || !exists || stored != lock {
			t.Fatalf("database failure discarded an active mutex: %+v", ret)
		}
	}
}
