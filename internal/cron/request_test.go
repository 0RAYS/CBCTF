package cron

import (
	"database/sql"
	"errors"
	"strings"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
)

func TestFailedRequestBatchIsRetainedWithStableIdentity(t *testing.T) {
	oldDB, oldLogger := db.CronDB, log.Logger
	previous := middleware.DrainRequestsPool()
	t.Cleanup(func() {
		db.CronDB, log.Logger = oldDB, oldLogger
		middleware.DrainRequestsPool()
		middleware.RestoreRequests(previous)
	})
	log.Init()
	database, err := gorm.Open(postgres.Open("host=127.0.0.1 dbname=test"), &gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := database.DB()
	defer func(pool *sql.DB) {
		_ = pool.Close()
	}(pool)
	db.CronDB = database
	failed := true
	var conflictSQL string
	_ = database.Callback().Create().After("gorm:create").Register("test:flush", func(tx *gorm.DB) {
		conflictSQL = tx.Statement.SQL.String()
		if failed {
			_ = tx.AddError(errors.New("database unavailable"))
		}
	})
	middleware.AppendRequest(model.Request{Path: "/first"})
	saveRequestLogTask()
	retained := middleware.DrainRequestsPool()
	if len(retained) != 1 || retained[0].RecordID == "" {
		t.Fatalf("failed batch lost: %+v", retained)
	}
	id := retained[0].RecordID
	middleware.RestoreRequests(retained)
	middleware.AppendRequest(model.Request{Path: "/second"})
	saveRequestLogTask()
	retained = middleware.DrainRequestsPool()
	if len(retained) != 2 || retained[0].RecordID != id {
		t.Fatalf("retry changed identity/order: %+v", retained)
	}
	middleware.RestoreRequests(retained)
	failed = false
	saveRequestLogTask()
	if len(middleware.DrainRequestsPool()) != 0 || !strings.Contains(conflictSQL, `ON CONFLICT ("record_id") DO NOTHING`) {
		t.Fatalf("successful flush not idempotent: %s", conflictSQL)
	}
}
