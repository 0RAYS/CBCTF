package resp

import (
	"github.com/gin-gonic/gin"

	"CBCTF/internal/model"
)

func GetCronJobResp(cronJob model.CronJob) gin.H {
	return gin.H{
		"id":           cronJob.ID,
		"name":         cronJob.Name,
		"description":  cronJob.Description,
		"schedule":     int64(cronJob.Schedule.Seconds()),
		"run_on_start": cronJob.RunOnStart,
		"success_last": cronJob.SuccessLast,
		"failure_last": cronJob.FailureLast,
		"success":      cronJob.Success,
		"failure":      cronJob.Failure,
	}
}
