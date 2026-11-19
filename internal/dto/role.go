package dto

type CreateRoleForm struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
}

type UpdateRoleForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
}

type AssignPermissionForm struct {
	PermissionID uint `form:"permission_id" json:"permission_id" binding:"required"`
}
