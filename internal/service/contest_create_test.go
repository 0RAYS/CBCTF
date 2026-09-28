package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/dto"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

type transactionRecorder struct {
	gorm.ConnPool
	begins, commits, rollbacks int
}

func (r *transactionRecorder) BeginTx(context.Context, *sql.TxOptions) (gorm.ConnPool, error) {
	r.begins++
	return &recordedTransaction{ConnPool: r.ConnPool, recorder: r}, nil
}

type recordedTransaction struct {
	gorm.ConnPool
	recorder *transactionRecorder
}

func (tx *recordedTransaction) Commit() error {
	tx.recorder.commits++
	return nil
}

func (tx *recordedTransaction) Rollback() error {
	tx.recorder.rollbacks++
	return nil
}

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
			transactions := &transactionRecorder{ConnPool: pool}
			tx.Statement.ConnPool = transactions
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
			contest, ret := CreateContest(tx, dto.CreateContestForm{Name: "Test", Hidden: hidden, Blood: blood})
			if !ret.OK || contest.Hidden != hidden || contest.Blood != blood || stored.Hidden != hidden || stored.Blood != blood {
				t.Fatalf("choices hidden=%v blood=%v were not retained: %+v, %+v", hidden, blood, contest, ret)
			}
			if transactions.begins != 1 || transactions.commits != 1 || transactions.rollbacks != 0 {
				t.Fatalf("expected one committed transaction: %+v", transactions)
			}
		}
	}
}

func TestMergedServiceTransactionsRollbackOnFailure(t *testing.T) {
	log.Init()
	for _, test := range []struct {
		name string
		call func(*gorm.DB) model.RetVal
	}{
		{"contest", func(tx *gorm.DB) model.RetVal {
			_, ret := CreateContest(tx, dto.CreateContestForm{Name: "Test"})
			return ret
		}},
		{"challenge", func(tx *gorm.DB) model.RetVal {
			_, ret := CreateChallenge(tx, dto.CreateChallengeForm{Name: "Test", Type: model.StaticChallengeType})
			return ret
		}},
		{"registration", func(tx *gorm.DB) model.RetVal {
			_, ret := RegisterUser(tx, dto.RegisterForm{Name: "Test", Password: "test-password"})
			return ret
		}},
		{"oauth", func(tx *gorm.DB) model.RetVal {
			_, ret := OauthLogin(tx, model.Oauth{IDClaim: "{id}"}, map[string]any{})
			return ret
		}},
		{"create team", func(tx *gorm.DB) model.RetVal {
			_, ret := CreateTeam(tx, model.Contest{Captcha: "expected"}, model.User{}, dto.CreateTeamForm{})
			return ret
		}},
		{"join team", func(tx *gorm.DB) model.RetVal {
			_, ret := JoinTeam(tx, model.Contest{}, model.User{}, dto.JoinTeamForm{Name: "missing"})
			return ret
		}},
		{"leave team", func(tx *gorm.DB) model.RetVal {
			return LeaveTeam(tx, model.Contest{}, model.Team{}, 1)
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			tx, err := gorm.Open(postgres.Open("host=127.0.0.1 user=test dbname=test"), &gorm.Config{
				DryRun: true, DisableAutomaticPing: true, SkipDefaultTransaction: true,
				Logger: logger.Default.LogMode(logger.Silent),
			})
			if err != nil {
				t.Fatal(err)
			}
			pool, _ := tx.DB()
			t.Cleanup(func() { _ = pool.Close() })
			transactions := &transactionRecorder{ConnPool: pool}
			tx.Statement.ConnPool = transactions
			if err := tx.Callback().Create().Before("gorm:create").Register("test:reject", func(query *gorm.DB) {
				query.AddError(errors.New("injected create failure"))
			}); err != nil {
				t.Fatal(err)
			}
			if ret := test.call(tx); ret.OK {
				t.Fatal("failure was swallowed")
			}
			if transactions.begins != 1 || transactions.commits != 0 || transactions.rollbacks != 1 {
				t.Fatalf("expected one rolled back transaction: %+v", transactions)
			}
		})
	}
}
