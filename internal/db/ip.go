package db

import (
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"context"
)

func RecordIP(ctx context.Context, ip model.IP) {
	if DB.WithContext(ctx).Create(&ip).Error != nil {
		log.Logger.Warningf("Failed to record IP %s", ip.IP)
	}
}

func CountIP(ctx context.Context) int64 {
	var count int64
	DB.WithContext(ctx).Model(&model.IP{}).Distinct("ip").Count(&count)
	return count
}
