package cron

import (
	"CBCTF/internal/log"
	"github.com/robfig/cron/v3"
	"time"
)

var c *cron.Cron

func executionTime(name string, task func()) func() {
	return func() {
		start := time.Now()
		log.Logger.Infof("%s processing time: %s", name, time.Since(start))
		task()
	}
}

func Init() {
	c = cron.New(cron.WithSeconds())
	CloseDockers(c)
	CloseUnCtrlDockers(c)
	UpdateUsageScore(c)
	UpdateUserRanking(c)
	UpdateTeamRanking(c)
	CheckCheat(c)
	CloseGenerator(c)
	PrepareGenerator(c)
	ClearUsageMutex(c)
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
