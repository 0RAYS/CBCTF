package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"github.com/robfig/cron/v3"
	"time"
)

// CloseDockers 关闭并删除超时 dockers
func CloseDockers(c *cron.Cron) {
	c.Schedule(cron.Every(1*time.Minute), cron.FuncJob(func() {
		dockers, ok, msg := db.GetDockers(db.DB, false)
		if !ok {
			log.Logger.Warningf("Failed to get dockers %s", msg)
			return
		}
		for _, docker := range dockers {
			if docker.Start.Add(docker.Duration).Before(time.Now()) {
				// 每次删除都作为一个单独的事务, 不回滚之前的删除
				tx := db.DB.Begin()
				if ok, msg = db.DeleteDocker(tx, docker); !ok {
					tx.Rollback()
					log.Logger.Warningf("Failed to delete docker %s", msg)
					continue
				}
				tx.Commit()
			}
		}
	}))
}
