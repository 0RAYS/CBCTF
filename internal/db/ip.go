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
