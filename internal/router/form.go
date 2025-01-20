package router

import "time"

type LoginForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegisterForm struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type GetModelsForm struct {
	Offset int `form:"offset"`
	Limit  int `form:"limit"`
}

type ChangePasswordForm struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

type UpdateSelfForm struct {
	Name    *string `json:"name"`
	Email   *string `json:"email"`
	Desc    *string `json:"desc"`
	Country *string `json:"country"`
	Website *string `json:"website"`
}

type UpdateAdminForm struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

type UpdateUserForm struct {
	Name     *string `json:"name"`
	Email    *string `json:"email"`
	Desc     *string `json:"desc"`
	Country  *string `json:"country"`
	Website  *string `json:"website"`
	Password *string `json:"password"`
	Hidden   *bool   `json:"hidden"`
	Banned   *bool   `json:"banned"`
	Verified *bool   `json:"verified"`
}

type UpdateTeamForm struct {
	Name      *string `json:"name"`
	Desc      *string `json:"desc"`
	Captcha   *string `json:"captcha"`
	CaptainID *uint   `json:"captain_id"`
}

type UpdateContestForm struct {
	Name     *string        `json:"name"`
	Desc     *string        `json:"desc"`
	Captcha  *string        `json:"captcha"`
	Avatar   *string        `json:"avatar"`
	Size     *int           `json:"size"`
	Start    *time.Time     `json:"start"`
	Duration *time.Duration `json:"duration"`
	Hidden   *bool          `json:"hidden"`
}

type AdminUpdateTeamForm struct {
	Name      *string `json:"name"`
	Desc      *string `json:"desc"`
	Hidden    *bool   `json:"hidden"`
	Banned    *bool   `json:"banned"`
	Captcha   *string `json:"captcha"`
	CaptainID *uint   `json:"captain_id"`
}

type DeleteSelfForm struct {
	Password string `json:"password" binding:"required"`
}

type DeleteUserForm struct {
	UserID uint `json:"user_id" binding:"required"`
}

type JoinTeamForm struct {
	Name    string `json:"name" binding:"required"`
	Captcha string `json:"captcha" binding:"required"`
}

type KickMemberForm struct {
	UserID uint `json:"user_id" binding:"required"`
}

type CreateTeamForm struct {
	Name string `json:"name" binding:"required"`
}

type CreateContestForm CreateTeamForm
