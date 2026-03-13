package db

import (
	"CBCTF/internal/i18n"
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

func (c *CronJobRepo) InitCronJob() model.RetVal {
	for _, cronJob := range model.CronJobs {
		res := c.DB.Model(&model.CronJob{}).FirstOrCreate(&cronJob, model.CronJob{Name: cronJob.Name})
		if res.Error != nil {
			return model.RetVal{Msg: i18n.Model.GetError, Attr: map[string]any{"Model": cronJob.ModelName(), "Error": res.Error.Error()}}
		}
	}
	return model.SuccessRetVal()
}
