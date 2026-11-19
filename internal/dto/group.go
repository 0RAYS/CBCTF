package dto

type CreateGroupForm struct {
	RoleID      uint   `form:"role_id" json:"role_id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
}

type UpdateGroupForm struct {
	RoleID      *uint   `form:"role_id" json:"role_id"`
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
}

type AssignUserGroupForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}
