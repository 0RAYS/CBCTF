package form

// LoginForm for user or admin login
type LoginForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// RegisterForm for user register
type RegisterForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

// CreateUserForm for create user
type CreateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Country  string `form:"country" json:"country"`
	Desc     string `form:"desc" json:"desc"`
	Hidden   bool   `form:"hidden" json:"hidden"`
	Verified bool   `form:"verified" json:"verified"`
	Banned   bool   `form:"banned" json:"banned"`
}

// UpdateSelfForm for user update info
type UpdateSelfForm struct {
	Name    *string `form:"name" json:"name"`
	Email   *string `form:"email" json:"email"`
	Desc    *string `form:"desc" json:"desc"`
	Country *string `form:"country" json:"country"`
}

// UpdateUserForm for admin update user info
type UpdateUserForm struct {
	Name     *string `form:"name" json:"name"`
	Email    *string `form:"name" json:"email"`
	Desc     *string `form:"desc" json:"desc"`
	Country  *string `form:"country" json:"country"`
	Password *string `form:"password" json:"password"`
	Hidden   *bool   `form:"hidden" json:"hidden"`
	Banned   *bool   `form:"banned" json:"banned"`
	Verified *bool   `form:"verified" json:"verified"`
}

// DeleteSelfForm for user delete self
type DeleteSelfForm struct {
	Password string `form:"password" json:"password" binding:"required"`
}

// DeleteUserForm for admin delete user
type DeleteUserForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}
