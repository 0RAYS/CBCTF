package cron

import (
	"CBCTF/internal/websocket"
	"github.com/robfig/cron/v3"
	"time"
)

func CheckWSConnection(c *cron.Cron) {
	function := exec("CheckWSConnection", func() {
		for id, conn := range websocket.AdminClients {
			if time.Now().Sub(conn.LastActive) > time.Minute {
				websocket.AdminClientsMu.Lock()
				delete(websocket.AdminClients, id)
				websocket.AdminClientsMu.Unlock()
			}
		}
		for id, conn := range websocket.UserClients {
			if time.Now().Sub(conn.LastActive) > time.Minute {
				websocket.UserClientsMu.Lock()
				delete(websocket.UserClients, id)
				websocket.UserClientsMu.Unlock()
			}
		}
	})
	c.Schedule(cron.Every(30*time.Second), cron.FuncJob(function))
}
