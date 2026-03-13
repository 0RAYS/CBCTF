package db

import (
	"CBCTF/internal/model"

	"gorm.io/gorm"
)

type CronJobRepo struct {
	BaseRepo[model.CronJob]
}

type UpdateCronJobOptions struct {
	Schedule *string
}

func (u UpdateCronJobOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Schedule != nil {
		options["schedule"] = *u.Schedule
	}
	return options
}

func InitCronJobRepo(tx *gorm.DB) *CronJobRepo {
	return &CronJobRepo{
		BaseRepo: BaseRepo[model.CronJob]{
			DB: tx,
		},
	}
}
