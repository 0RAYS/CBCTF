package cron

import (
	"errors"
	"testing"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func TestFlagScoreReadFailureMarksCronFailed(t *testing.T) {
	old, oldLogger := db.CronDB, log.Logger
	t.Cleanup(func() { db.CronDB, log.Logger = old, oldLogger })
	log.Init()
	database, err := gorm.Open(postgres.Open("host=127.0.0.1 dbname=test"), &gorm.Config{DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true, Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	pool, _ := database.DB()
	defer pool.Close()
	db.CronDB = database
	_ = database.Callback().Query().After("gorm:query").Register("test:score_reads", func(tx *gorm.DB) {
		switch dst := tx.Statement.Dest.(type) {
		case *model.CronJob:
			dst.ID = 1
			dst.Schedule = time.Minute
			tx.RowsAffected = 1
		case *[]model.Contest:
			*dst = []model.Contest{{ID: 1, Start: time.Now(), Duration: time.Hour}}
			tx.RowsAffected = 1
		default:
			_ = tx.AddError(errors.New("score source unavailable"))
		}
	})
	if ret := updateFlagScoreTask(); ret.OK {
		t.Fatal("score lookup failure reported a successful cron run")
	}
}
