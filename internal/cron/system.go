package cron

import (
	"CBCTF/internal/log"
	"CBCTF/internal/redis"
)

// collectSystemMetricsTask 收集系统运行状态
func collectSystemMetricsTask() {
	metrics, err := redis.CollectMetrics()
	if err != nil {
		log.Logger.Warningf("Failed to collect system metrics: %s", err)
		return
	}
	if err = redis.SaveMetrics(metrics); err != nil {
		log.Logger.Warningf("Failed to save system metrics: %s", err)
		return
	}
}
