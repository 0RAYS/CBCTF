package dto

type UpdatePermissionForm struct {
	Description *string `form:"description" json:"description"`
}
