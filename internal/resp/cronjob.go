package resp

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

func GetCronJobResp(cronJob model.CronJob) gin.H {
	return gin.H{
		"id":          cronJob.ID,
		"name":        cronJob.Name,
		"description": cronJob.Description,
		"schedule":    cronJob.Schedule,
		"status":      cronJob.Status,
	}
}
