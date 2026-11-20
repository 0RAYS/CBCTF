package db

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

type CronJobRepo struct {
	BaseRepo[model.CronJob]
}

type UpdateCronJobOptions struct {
	Schedule    *time.Duration
	RunOnStart  *bool
	Success     *int64
	SuccessLast *time.Time
	Failure     *int64
	FailureLast *time.Time
}

func (u UpdateCronJobOptions) Convert2Map() map[string]any {
	options := make(map[string]any)
	if u.Schedule != nil {
		options["schedule"] = *u.Schedule
	}
	if u.RunOnStart != nil {
		options["run_on_start"] = *u.RunOnStart
	}
	if u.Success != nil {
		options["success"] = *u.Success
	}
	if u.SuccessLast != nil {
		options["success_last"] = *u.SuccessLast
	}
	if u.Failure != nil {
		options["failure"] = *u.Failure
	}
	if u.FailureLast != nil {
		options["failure_last"] = *u.FailureLast
	}
	return options
}

type DiffUpdateCronJobOptions struct {
	Success int64
	Failure int64
}

func (d DiffUpdateCronJobOptions) Convert2Expr() map[string]any {
	options := make(map[string]any)
	if d.Success != 0 {
		options["success"] = gorm.Expr("success + ?", d.Success)
	}
	if d.Failure != 0 {
		options["failure"] = gorm.Expr("failure + ?", d.Failure)
	}
	return options
}

func InitCronJobRepo(tx *gorm.DB) *CronJobRepo {
	return &CronJobRepo{
		DB: tx,
	}
}

func (c *CronJobRepo) InitCronJob() model.RetVal {
	for _, cronJob := range model.CronJobs {
		cronJob.DefaultSchedule = cronJob.Schedule
		target := fmt.Sprintf("name=%s", cronJob.Name)
		var existing model.CronJob
		res := c.DB.Model(&model.CronJob{}).Where("name = ?", cronJob.Name).First(&existing)
		if res.Error != nil {
			if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return model.FailRetVal(i18n.Model.CronJob.GetError, "CronJob", target, res.Error)
			}
			if err := c.DB.Model(&model.CronJob{}).Create(&cronJob).Error; err != nil {
				return model.FailRetVal(i18n.Model.CronJob.GetError, "CronJob", target, err)
			}
			continue
		}

		updates := map[string]any{
			"description":      cronJob.Description,
			"default_schedule": cronJob.DefaultSchedule,
		}
		if existing.Schedule <= 0 || existing.Schedule == existing.DefaultSchedule {
			updates["schedule"] = cronJob.Schedule
		}
		if err := c.DB.Model(&existing).Updates(updates).Error; err != nil {
			return model.FailRetVal(i18n.Model.CronJob.GetError, "CronJob", target, err)
		}
	}
	return model.SuccessRetVal()
}

func (c *CronJobRepo) UpdateStatus(id uint, success bool, last time.Time) model.RetVal {
	return updateRunStatus(&c.BaseRepo, id, success, last)
}
