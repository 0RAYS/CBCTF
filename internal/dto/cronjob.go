package dto

type UpdateCronJobForm struct {
	Schedule *int64 `form:"schedule" json:"schedule" binding:"omitempty,gte=1"`
}
