package dto

// CreateRoleForm 创建角色表单
type CreateRoleForm struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
}

// UpdateRoleForm 更新角色表单
type UpdateRoleForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
}

// AssignPermissionForm 分配/撤销权限表单
type AssignPermissionForm struct {
	PermissionID uint `form:"permission_id" json:"permission_id" binding:"required"`
}
