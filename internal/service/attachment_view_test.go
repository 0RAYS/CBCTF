package service

import (
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"CBCTF/internal/config"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
)

func TestChallengeListAndStatusUseDownloadNameInsteadOfServerPath(t *testing.T) {
	previous, previousLog := config.Env, log.Logger
	config.Env = &config.Config{Path: t.TempDir()}
	log.Init()
	t.Cleanup(func() { config.Env, log.Logger = previous, previousLog })
	// Only file presentation is exercised; database queries run without a connection.
	tx, err := gorm.Open(postgres.Open("host=127.0.0.1 user=test dbname=test"), &gorm.Config{
		DryRun: true, DisableAutomaticPing: true, Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}
	pool, err := tx.DB()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = pool.Close() })
	challenge := model.Challenge{ID: 3, Type: model.StaticChallengeType}
	if err := os.MkdirAll(filepath.Dir(challenge.StaticPath()), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(challenge.StaticPath(), []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	contestChallenge := model.ContestChallenge{Challenge: challenge}
	list := BuildContestChallengeRuntimeView(tx, model.Team{}, contestChallenge)
	status, ret := GetContestChallengeStatus(tx, model.Team{}, challenge, contestChallenge)
	if !ret.OK {
		t.Fatal(ret)
	}
	if list.FileName != model.AttachmentFileName || status.FileName != list.FileName {
		t.Fatalf("list file=%q, status file=%q", list.FileName, status.FileName)
	}
}
