package dto

type UpdateCronJobForm struct {
	Schedule *string `form:"schedule" json:"schedule" binding:"omitempty,max=255"`
}
