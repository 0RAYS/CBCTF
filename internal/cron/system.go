package cron

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/redis"
)

// collectSystemMetricsTask 收集系统运行状态
func collectSystemMetricsTask() model.RetVal {
	metrics, err := redis.CollectMetrics()
	if err != nil {
		log.Logger.Warningf("Failed to collect system metrics: %s", err)
		return model.RetVal{Msg: err.Error()}
	}
	if err = redis.SaveMetrics(metrics); err != nil {
		log.Logger.Warningf("Failed to save system metrics: %s", err)
		return model.RetVal{Msg: err.Error()}
	}
	return model.SuccessRetVal()
}
