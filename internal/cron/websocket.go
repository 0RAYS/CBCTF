package cron

import (
	"CBCTF/internal/websocket"
	"time"

	"github.com/robfig/cron/v3"
)

func CheckWSConnection(c *cron.Cron) {
	c.Schedule(cron.Every(5*time.Second), cron.FuncJob(exec("CheckWSConnection", func() {
		for id, conn := range websocket.AdminClients {
			if conn.LastActive.Add(10 * time.Second).Before(time.Now()) {
				websocket.AdminClientsMu.Lock()
				delete(websocket.AdminClients, id)
				websocket.AdminClientsMu.Unlock()
			}
		}
		for id, conn := range websocket.UserClients {
			if conn.LastActive.Add(10 * time.Second).Before(time.Now()) {
				websocket.UserClientsMu.Lock()
				delete(websocket.UserClients, id)
				websocket.UserClientsMu.Unlock()
			}
		}
	})))
}
