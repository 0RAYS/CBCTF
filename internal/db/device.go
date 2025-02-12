package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/gorm"
	"strings"
)

// RecordDevice 记录设备 ID
func RecordDevice(tx *gorm.DB, device model.Device) {
	if err := tx.Model(&model.Device{}).Create(&device).Error; err != nil && !strings.Contains(err.Error(), "Error 1062") {
		log.Logger.Warningf("Failed to record device %d-%s: %s", device.UserID, device.Magic, err)
	}
}
