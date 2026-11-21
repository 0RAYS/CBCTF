package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"CBCTF/internal/dto"
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type failedCommitPool struct{ gorm.ConnPool }

func (p *failedCommitPool) BeginTx(context.Context, *sql.TxOptions) (gorm.ConnPool, error) {
	return &failedCommitTx{ConnPool: p.ConnPool}, nil
}

type failedCommitTx struct{ gorm.ConnPool }

func (*failedCommitTx) Commit() error   { return errors.New("commit failed") }
func (*failedCommitTx) Rollback() error { return nil }

func TestContestImportDoesNotPublishUncommittedSuccess(t *testing.T) {
	database := queryFailureDB(t)
	pool, _ := database.DB()
	database.Statement.ConnPool = &failedCommitPool{ConnPool: pool}
	created := false
	_ = database.Callback().Create().After("gorm:create").Register("test:create_import", func(tx *gorm.DB) {
		if value, ok := tx.Statement.Dest.(*model.Challenge); ok {
			value.ID = 1
			created = true
			tx.RowsAffected = 1
		}
	})
	_ = database.Callback().Query().After("gorm:query").Register("test:import_rows", func(tx *gorm.DB) {
		dst, ok := tx.Statement.Dest.(*model.Challenge)
		if !ok {
			return
		}
		for _, value := range tx.Statement.Vars {
			switch typed := value.(type) {
			case string:
				if typed != "one" {
					continue
				}
			case uint:
				if typed != 1 {
					continue
				}
			default:
				continue
			}
			*dst = model.Challenge{ID: 1, RandID: "one", Type: model.StaticChallengeType}
			tx.RowsAffected = 1
			return
		}
	})
	items, failed, ret := CopyChallengesToContest(database, model.Contest{ID: 1}, dto.CopyChallengesToContestForm{ChallengeIDs: []string{"one", "two"}})
	batch, ok := ret.Data.(*model.BatchResult)
	if ret.OK || !created || len(items) != 0 || len(failed) != 1 || !ok || batch.NotAttempted != 1 {
		t.Fatalf("commit failure reported success: %+v %v %+v", items, failed, ret)
	}
}
