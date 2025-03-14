package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/gorm"
	"strings"
)

// RecordDevice 记录设备 ID
func RecordDevice(tx *gorm.DB, device model.Device) {
	var tmp model.Device
	res := tx.Model(&model.Device{}).
		Where("user_id = ? AND magic = ?", device.UserID, device.Magic).Find(&tmp).Limit(1)
	if res.RowsAffected == 1 {
		tmp.Count = tmp.Count + 1
		if err := tx.Save(&tmp).Error; err != nil {
			log.Logger.Warningf("Failed to record device %d-%s: %s", device.UserID, device.Magic, err)
		}
		return
	}
	device.Count = 1
	if err := tx.Model(&model.Device{}).Create(&device).Error; err != nil && !strings.Contains(err.Error(), "Error 1062") {
		log.Logger.Warningf("Failed to record device %d-%s: %s", device.UserID, device.Magic, err)
	}
}
