package cron

import (
	"CBCTF/internel/log"
	"github.com/robfig/cron/v3"
	"time"
)

var c *cron.Cron

func exec(name string, task func()) func() {
	return func() {
		start := time.Now()
		task()
		log.Logger.Infof("%s processing time: %s", name, time.Since(start))
	}
}

func Init() {
	c = cron.New(cron.WithSeconds())
	StopTimeoutVictims(c)
	StopUnCtrlPods(c)
	ClearUnCtrlResource(c)
	UpdateFlagScore(c)
	UpdateUserRanking(c)
	UpdateTeamRanking(c)
	ResetGenerator(c)
	ClearUsageMutex(c)
	//CheckCheat(c)
}

func Start() {
	log.Logger.Info("Cron started")
	c.Start()
}

func Stop() {
	if c != nil {
		c.Stop()
		log.Logger.Info("Cron stopped")
	}
}
