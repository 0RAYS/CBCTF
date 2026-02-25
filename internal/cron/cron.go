package cron

import (
	"CBCTF/internal/log"
	"CBCTF/internal/prometheus"
	"time"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

func exec(name string, task func() error) func() {
	return func() {
		start := time.Now()
		err := task()
		duration := time.Since(start).Seconds()
		prometheus.RecordCronJob(name, duration, err == nil)
		if err != nil {
			log.Logger.Warningf("%s failed: %s, processing time: %s", name, err, time.Duration(duration*float64(time.Second)))
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
	checkWSConnection(c)
	closeTimeoutVictims(c)
	closeUnCtrlVictims(c)
	ClearEmptyTeam(c)
	updateFlagScore(c)
	updateUserRanking(c)
	updateTeamRanking(c)
	stopUnCtrlGenerator(c)
	prepareGenerator(c)
	clearSubmissionMutex(c)
	checkCheat(c)
	clearCheatMutex(c)
	clearJoinTeamMutes(c)
	updateJWTSecret(c)
	c.Start()
}

func Stop() {
	c.Stop()
}
