package cron

import (
	"CBCTF/internal/config"
	"CBCTF/internal/utils"
	"time"

	"github.com/robfig/cron/v3"
)

func updateJWTSecret(c *cron.Cron) {
	c.Schedule(cron.Every(time.Hour*2), cron.FuncJob(exec("UpdateJWTSecret", func() error {
		if !config.Env.Gin.JWT.Static || config.Env.Gin.JWT.Secret == "" {
			config.Env.Gin.JWT.Secret = utils.UUID()
			config.Env.Gin.JWT.Static = false
		}
		return nil
	})))
}
