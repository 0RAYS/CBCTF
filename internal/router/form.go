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

type UserUpdateForm struct {
	Name    *string `json:"name"`
	Email   *string `json:"email"`
	Desc    *string `json:"desc"`
	Country *string `json:"country"`
	Website *string `json:"website"`
}

type AdminUpdateForm struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type GetContestsForm struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type GetUsersForm GetContestsForm

type GetTeamsForm struct {
	Offset    int  `form:"offset"`
	Limit     int  `form:"limit"`
	ContestID uint `form:"contest_id"`
}
