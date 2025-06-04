package form

// CreateAdminForm for create admin
type CreateAdminForm RegisterForm

// UpdateAdminForm for admin update info
type UpdateAdminForm struct {
	Name  *string `form:"name" json:"name"`
	Email *string `form:"email" json:"email"`
}
