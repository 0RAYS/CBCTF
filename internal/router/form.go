package router

type LoginForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type ChangePasswordForm struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type UpdateForm struct {
	Name    *string `json:"name"`
	Email   *string `json:"email"`
	Desc    *string `json:"desc"`
	Country *string `json:"country"`
	Website *string `json:"website"`
}
