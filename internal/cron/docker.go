package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"context"
	"github.com/robfig/cron/v3"
	"time"
)

func CloseDockers(c *cron.Cron) {
	c.Schedule(cron.Every(1*time.Minute), cron.FuncJob(func() {
		dockers, ok, msg := db.GetDockers(context.Background())
		if !ok {
			log.Logger.Warningf("Failed to get dockers %s", msg)
			return
		}
		for _, docker := range dockers {
			if docker.Start.Add(docker.Duration).Before(time.Now()) {
				_, _ = db.DeleteDocker(context.Background(), docker)
			}
		}
	}))
}
