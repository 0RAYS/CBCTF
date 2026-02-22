package dto

// UpdatePermissionForm 更新权限表单
type UpdatePermissionForm struct {
	Description *string `form:"description" json:"description"`
}
