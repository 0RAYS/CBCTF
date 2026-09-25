package service

import (
	"time"

	"gorm.io/gorm"

	"CBCTF/internal/db"
	"CBCTF/internal/dto"
	"CBCTF/internal/model"
)

func ListCronJobs(tx *gorm.DB, form dto.ListModelsForm) ([]model.CronJob, int64, model.RetVal) {
	return db.InitCronJobRepo(tx).List(form.Limit, form.Offset)
}

func UpdateCronJob(tx *gorm.DB, cronJob model.CronJob, form dto.UpdateCronJobForm) (model.CronJob, model.RetVal) {
	var schedule *time.Duration
	if form.Schedule != nil {
		schedule = new(time.Duration(*form.Schedule) * time.Second)
	}
	if ret := db.InitCronJobRepo(tx).Update(cronJob.ID, db.UpdateCronJobOptions{
		Schedule: schedule,
	}); !ret.OK {
		return model.CronJob{}, ret
	}
	newCronJob, ret := db.InitCronJobRepo(tx).GetByID(cronJob.ID)
	if !ret.OK {
		return model.CronJob{}, ret
	}
	return newCronJob, model.SuccessRetVal()
}
