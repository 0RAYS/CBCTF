package dto

import "time"

type UpdateCronJobForm struct {
	Schedule *time.Duration `form:"schedule" json:"schedule" binding:"omitempty,gte=1000000000"`
}
