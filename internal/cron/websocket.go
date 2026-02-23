package cron

import (
	"CBCTF/internal/websocket"
	"time"

	"github.com/robfig/cron/v3"
)

// checkWSConnection 清理长时未通信的连接
func checkWSConnection(c *cron.Cron) {
	c.Schedule(cron.Every(5*time.Second), cron.FuncJob(func() {
		websocket.UserClientsMu.Lock()
		for id, conn := range websocket.UserClients {
			if conn.LastActive.Add(10 * time.Second).Before(time.Now()) {
				delete(websocket.UserClients, id)
			}
		}
		websocket.UserClientsMu.Unlock()
	}))
}
