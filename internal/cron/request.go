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

	if ret := db.InitRequestRepo(db.CronDB).Create(requests...); !ret.OK {
		middleware.RestoreRequests(requests)
		log.Logger.Warningf("Request log batch retained for retry: count=%d reason=%s", len(requests), ret.Msg)
	}
}

func saveTaskLogTask() {
	records := task.DrainTaskRecordPool()
	if len(records) == 0 {
		return
	}

	if ret := db.InitTaskRepo(db.TaskDB).CreateBatch(records...); !ret.OK {
		task.RestoreTaskRecords(records)
		log.Logger.Warningf("Failed to save task history batch: %s", ret.Msg)
	}
}
