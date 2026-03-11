package cron

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/prometheus"
	"time"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func exec(name string, task func() model.RetVal) func() {
	return func() {
		start := time.Now()
		ret := task()
		duration := time.Since(start).Seconds()
		prometheus.RecordCronJob(name, duration, ret.OK)
		if !ret.OK {
			log.Logger.Warningf("%s failed: %s, processing time: %s", name, ret.Msg, time.Duration(duration*float64(time.Second)))
		} else if duration > time.Second.Seconds() {
			log.Logger.Warningf("%s processing time: %s", name, time.Duration(duration*float64(time.Second)))
		} else {
			log.Logger.Debugf("%s processing time: %s", name, time.Duration(duration*float64(time.Second)))
		}
	}
}

func Init() {
	c = cron.New(cron.WithSeconds())
}

func Start() {
	log.Logger.Info("Cron started")
	collectSystemMetrics(c)
	closeTimeoutVictims(c)
	closeUnCtrlVictims(c)
	ClearEmptyTeam(c)
	updateFlagScore(c)
	updateUserRanking(c)
	updateTeamRanking(c)
	stopUnCtrlGenerator(c)
	clearSubmissionMutex(c)
	checkCheat(c)
	clearCheatMutex(c)
	clearJoinTeamMutes(c)
	c.Start()
}

func Stop() {
	c.Stop()
}
