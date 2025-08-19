package cron

import (
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
	"time"

	"github.com/robfig/cron/v3"
)

// collectSystemMetrics 收集系统运行状态
func collectSystemMetrics(c *cron.Cron) {
	c.Schedule(cron.Every(time.Second), cron.FuncJob(func() {
		metrics, err := redis.CollectMetrics()
		if err != nil {
			log.Logger.Warningf("Failed to collect system metrics: %s", err)
			return
		}
		if err = redis.SaveMetrics(metrics); err != nil {
			log.Logger.Warningf("Failed to save system metrics: %s", err)
		}
	}))
}
