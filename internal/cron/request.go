package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
)

func saveRequestLogTask() {
	requests := middleware.DrainRequestsPool()
	if len(requests) == 0 {
		return
	}

	db.InitRequestRepo(db.DB).Insert(requests...)
}
