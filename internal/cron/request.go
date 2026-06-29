package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/log"
	"CBCTF/internal/middleware"
	"CBCTF/internal/task"
)

func saveRequestLogTask() {
	requests := middleware.DrainRequestsPool()
	if len(requests) == 0 {
		return
	}

	db.InitRequestRepo(db.CronDB).Create(requests...)
}

func saveTaskLogTask() {
	records := task.DrainTaskRecordPool()
	if len(records) == 0 {
		return
	}

	if ret := db.InitTaskRepo(db.TaskDB).CreateBatch(records...); !ret.OK {
		log.Logger.Warningf("Failed to save task history batch: %s", ret.Msg)
	}
}
