package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
)

func saveRequestLogTask() {
	middleware.RequestsMutex.RLock()
	if len(middleware.RequestsPool) == 0 {
		middleware.RequestsMutex.RUnlock()
		return
	}
	db.InitRequestRepo(db.DB).Insert(middleware.RequestsPool...)
	middleware.RequestsMutex.RUnlock()
	middleware.RequestsMutex.Lock()
	middleware.RequestsPool = make([]model.Request, 0)
	middleware.RequestsMutex.Unlock()
}
