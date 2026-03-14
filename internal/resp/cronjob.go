package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetCronJobResp(cronJob model.CronJob) gin.H {
	return gin.H{
		"id":           cronJob.ID,
		"name":         cronJob.Name,
		"description":  cronJob.Description,
		"schedule":     cronJob.Schedule.String(),
		"schedule_ns":  int64(cronJob.Schedule),
		"success_last": cronJob.SuccessLast,
		"failure_last": cronJob.FailureLast,
		"success":      cronJob.Success,
		"failure":      cronJob.Failure,
	}
}
