package db

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
	"errors"
	"time"

	"gorm.io/gorm"
)

type CronJobRepo struct {
	BaseRepo[model.CronJob]
}

type UpdateCronJobOptions struct {
	Schedule    *time.Duration
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
		BaseRepo: BaseRepo[model.CronJob]{
			DB: tx,
		},
	}
}

func (c *CronJobRepo) InitCronJob() model.RetVal {
	for _, cronJob := range model.CronJobs {
		cronJob.DefaultSchedule = cronJob.Schedule
		var existing model.CronJob
		res := c.DB.Model(&model.CronJob{}).Where("name = ?", cronJob.Name).First(&existing)
		if res.Error != nil {
			if !errors.Is(res.Error, gorm.ErrRecordNotFound) {
				return model.RetVal{Msg: i18n.Model.CronJob.GetError, Attr: map[string]any{"Error": res.Error.Error()}}
			}
			if err := c.DB.Model(&model.CronJob{}).Create(&cronJob).Error; err != nil {
				return model.RetVal{Msg: i18n.Model.CronJob.GetError, Attr: map[string]any{"Error": err.Error()}}
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
			return model.RetVal{Msg: i18n.Model.CronJob.GetError, Attr: map[string]any{"Error": err.Error()}}
		}
	}
	return model.SuccessRetVal()
}

func (c *CronJobRepo) UpdateStatus(id uint, success bool, last time.Time) model.RetVal {
	var diffOptions DiffUpdateCronJobOptions
	var options UpdateCronJobOptions
	if success {
		diffOptions = DiffUpdateCronJobOptions{
			Success: 1,
		}
		options = UpdateCronJobOptions{
			SuccessLast: &last,
		}
	} else {
		diffOptions = DiffUpdateCronJobOptions{
			Failure: 1,
		}
		options = UpdateCronJobOptions{
			FailureLast: &last,
		}
	}
	if ret := c.DiffUpdate(id, diffOptions); !ret.OK {
		return ret
	}
	return c.Update(id, options)
}
