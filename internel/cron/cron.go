package cron

import (
	"CBCTF/internel/log"
	"github.com/robfig/cron/v3"
	"time"
)

var Cron *cron.Cron

func exec(name string, task func()) func() {
	return func() {
		start := time.Now()
		task()
		log.Logger.Infof("%s processing time: %s", name, time.Since(start))
	}
}

func Init() {
	Cron = cron.New(cron.WithSeconds())
	CloseTimeoutVictims(Cron)
	CloseUnCtrlVictims(Cron)
	ClearUnCtrlResource(Cron)
	UpdateFlagScore(Cron)
	UpdateUserRanking(Cron)
	UpdateTeamRanking(Cron)
	StopUnCtrlGenerator(Cron)
	PrepareGenerator(Cron)
	ClearContestChallengeMutex(Cron)
	CheckCheat(Cron)
}

func Start() {
	log.Logger.Info("Cron started")
	Cron.Start()
}
