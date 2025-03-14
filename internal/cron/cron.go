package cron

import (
	"CBCTF/internal/log"
	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func Init() {
	c = cron.New(cron.WithSeconds())
	CloseDockers(c)
	CloseUnCtrlDockers(c)
	UpdateUserScore(c)
	// UpdateGlobalRanking UpdateTeamRank 顺序不可调换
	UpdateGlobalRanking(c)
	UpdateTeamRank(c)
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
