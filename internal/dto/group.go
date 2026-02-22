package dto

// CreateGroupForm 创建分组表单
type CreateGroupForm struct {
	RoleID      uint   `form:"role_id" json:"role_id" binding:"required"`
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
}

// UpdateGroupForm 更新分组表单
type UpdateGroupForm struct {
	RoleID      *uint   `form:"role_id" json:"role_id"`
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
}

// AssignUserGroupForm 分配/移除用户分组表单
type AssignUserGroupForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}
