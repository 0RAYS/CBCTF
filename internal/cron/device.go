package cron

import (
	"CBCTF/internal/db"
	"CBCTF/internal/middleware"
	"CBCTF/internal/model"
)

func saveRequestDeviceTask() {
	devices := middleware.DrainRequestDevicePool()
	if len(devices) == 0 {
		return
	}

	deviceRepo := db.InitDeviceRepo(db.DB)
	for _, device := range devices {
		deviceRepo.RecordDevice(model.Device{
			UserID: device.UserID,
			Magic:  device.Magic,
		}, device.Count)
	}
}
