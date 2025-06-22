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
		duration := time.Since(start)
		if duration > time.Second {
			log.Logger.Warningf("%s processing time: %s", name, duration)
		} else {
			log.Logger.Debugf("%s processing time: %s", name, duration)
		}
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
