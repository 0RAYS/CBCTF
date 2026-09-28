package dto

type UpdateCronJobForm struct {
	Schedule   *int64 `form:"schedule" json:"schedule" binding:"omitempty,gte=1"`
	RunOnStart *bool  `form:"run_on_start" json:"run_on_start"`
}
