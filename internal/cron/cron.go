package cron

import (
	"CBCTF/internal/log"
	"time"

	"github.com/robfig/cron/v3"
)

var Cron *cron.Cron

func exec(name string, task func()) func() {
	return func() {
		start := time.Now()
		task()
		if duration := time.Since(start); duration > time.Second {
			log.Logger.Warningf("%s processing time: %s", name, duration)
		} else {
			log.Logger.Debugf("%s processing time: %s", name, duration)
		}
	}
}

func Init() {
	Cron = cron.New(cron.WithSeconds())
	CheckWSConnection(Cron)
	CloseTimeoutVictims(Cron)
	CloseUnCtrlVictims(Cron)
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
