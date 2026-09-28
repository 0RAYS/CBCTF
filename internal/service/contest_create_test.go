package service

import (
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/dto"
	"CBCTF/internal/model"
)

// Use GORM's real PostgreSQL statement generation, with callback-backed rows,
// to exercise default:true on INSERT and explicit false on UPDATE without a DB.
func TestCreateContestPreservesBooleanChoices(t *testing.T) {
	for _, hidden := range []bool{false, true} {
		for _, blood := range []bool{false, true} {
			tx, err := gorm.Open(postgres.Open("host=127.0.0.1 user=test dbname=test"), &gorm.Config{
				DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true,
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatal(err)
			}
			pool, _ := tx.DB()
			t.Cleanup(func() { _ = pool.Close() })
			var stored *model.Contest
			if err := tx.Callback().Create().After("gorm:create").Register("test:store", func(query *gorm.DB) {
				if contest, ok := query.Statement.Dest.(*model.Contest); ok {
					contest.ID = 1
					copy := *contest
					stored = &copy
					query.RowsAffected = 1
				}
			}); err != nil {
				t.Fatal(err)
			}
			if err := tx.Callback().Query().After("gorm:query").Register("test:load", func(query *gorm.DB) {
				if contest, ok := query.Statement.Dest.(*model.Contest); ok && stored != nil {
					*contest = *stored
					query.RowsAffected = 1
				}
			}); err != nil {
				t.Fatal(err)
			}
			if err := tx.Callback().Update().After("gorm:update").Register("test:update", func(query *gorm.DB) {
				changes := query.Statement.Dest.(map[string]any)
				stored.Hidden = changes["hidden"].(bool)
				stored.Blood = changes["blood"].(bool)
				query.RowsAffected = 1
			}); err != nil {
				t.Fatal(err)
			}
			contest, ret := createContest(tx, dto.CreateContestForm{Name: "Test", Hidden: hidden, Blood: blood})
			if !ret.OK || contest.Hidden != hidden || contest.Blood != blood || stored.Hidden != hidden || stored.Blood != blood {
				t.Fatalf("choices hidden=%v blood=%v were not retained: %+v, %+v", hidden, blood, contest, ret)
			}
		}
	}
}
