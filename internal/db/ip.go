package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"gorm.io/gorm"
)

// RecordIP 记录请求的 IP
func RecordIP(tx *gorm.DB, ip model.IP) {
	if err := tx.Model(&model.IP{}).Create(&ip).Error; err != nil {
		log.Logger.Warningf("Failed to record IP %s: %s", ip.IP, err)
	}
}

// CountIP 获取 IP 数量
func CountIP(tx *gorm.DB) int64 {
	var count int64
	tx.Model(&model.IP{}).Distinct("ip").Count(&count)
	return count
}
