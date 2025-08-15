package cron

import (
	"CBCTF/internal/log"
	"time"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

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
	c = cron.New(cron.WithSeconds())
}

func Start() {
	log.Logger.Info("Cron started")
	CollectSystemMetrics(c)
	CheckWSConnection(c)
	CloseTimeoutVictims(c)
	CloseUnCtrlVictims(c)
	UpdateFlagScore(c)
	UpdateUserRanking(c)
	UpdateTeamRanking(c)
	StopUnCtrlGenerator(c)
	PrepareGenerator(c)
	ClearContestChallengeMutex(c)
	CheckCheat(c)
	c.Start()
}
