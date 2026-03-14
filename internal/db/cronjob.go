package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"time"

	"gorm.io/gorm"
)

type CronJobRepo struct {
	BaseRepo[model.CronJob]
}

type UpdateCronJobOptions struct {
	Schedule        *time.Duration
	SuccessLast     *time.Time
	FailureLast     *time.Time
	SuccessCount    *uint
	FailureCount    *uint
	IncreaseSuccess bool
	IncreaseFailure bool
}

func (u UpdateCronJobOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Schedule != nil {
		options["schedule"] = *u.Schedule
	}
	if u.SuccessLast != nil {
		options["success_last"] = *u.SuccessLast
	}
	if u.FailureLast != nil {
		options["failure_last"] = *u.FailureLast
	}
	if u.SuccessCount != nil {
		options["success_count"] = *u.SuccessCount
	}
	if u.FailureCount != nil {
		options["failure_count"] = *u.FailureCount
	}
	if u.IncreaseSuccess {
		options["success_count"] = gorm.Expr("success_count + ?", 1)
	}
	if u.IncreaseFailure {
		options["failure_count"] = gorm.Expr("failure_count + ?", 1)
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

func (c *CronJobRepo) UpdateByName(name string, options UpdateCronJobOptions) model.RetVal {
	cronJob, ret := c.GetByUniqueKey("name", name)
	if !ret.OK {
		return ret
	}
	return c.Update(cronJob.ID, options)
}
