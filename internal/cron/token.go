package cron

import (
	"CBCTF/internal/utils"
	"time"

	"github.com/robfig/cron/v3"
)

func updateJWTSecret(c *cron.Cron) {
	c.Schedule(cron.Every(time.Hour*12), cron.FuncJob(exec("UpdateJWTSecret", func() {
		utils.JWTSecret = utils.UUID()
	})))
}
