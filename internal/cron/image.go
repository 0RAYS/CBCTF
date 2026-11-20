package cron

import (
	"time"

	"CBCTF/internal/db"
	"CBCTF/internal/model"
	"CBCTF/internal/service"
)

// detectImagesTask detects the images of the whole challenge library and of the
// unfinished contests. The run merges into the common detection cache, so the
// system-wide view and the contest views share the results.
func detectImagesTask() model.RetVal {
	images, ret := service.DetectionImages(db.CronDB, time.Now())
	if !ret.OK {
		return ret
	}
	return service.DetectImages(db.CronDB.Statement.Context, images)
}
