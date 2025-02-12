package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"errors"
	"gorm.io/gorm"
)

// RecordDevice 记录设备 ID
func RecordDevice(tx *gorm.DB, device model.Device) {
	if err := tx.Model(&model.Device{}).Create(&device).Error; err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
		log.Logger.Warningf("Failed to record IP %d-%s", device.UserID, device.Magic)
	}
}
