package form

// GetModelsForm for get models list
type GetModelsForm struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

// ChangePasswordForm for user or admin change password
type ChangePasswordForm struct {
	OldPassword string `form:"old" json:"old" binding:"required"`
	NewPassword string `form:"new" json:"new" binding:"required"`
}
