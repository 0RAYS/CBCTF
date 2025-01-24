package router

import "time"

type LoginForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type RegisterForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
}

type GetModelsForm struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

type ChangePasswordForm struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required"`
}

type UpdateSelfForm struct {
	Name    *string `form:"name" json:"name"`
	Email   *string `form:"email" json:"email"`
	Desc    *string `form:"desc" json:"desc"`
	Country *string `form:"country" json:"country"`
}

type UpdateAdminForm struct {
	Name  *string `form:"name" json:"name"`
	Email *string `form:"email" json:"email"`
}

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

type UpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	Captcha   *string `form:"captcha" json:"captcha"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
}

type UpdateContestForm struct {
	Name     *string        `form:"name" json:"name"`
	Desc     *string        `form:"desc" json:"desc"`
	Captcha  *string        `form:"captcha" json:"captcha"`
	Avatar   *string        `form:"avatar" json:"avatar"`
	Size     *int           `form:"start" json:"size"`
	Start    *time.Time     `form:"start" json:"start"`
	Duration *time.Duration `form:"duration" json:"duration"`
	Hidden   *bool          `form:"hidden" json:"hidden"`
}

type AdminUpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	Hidden    *bool   `form:"hidden" json:"hidden"`
	Banned    *bool   `form:"banned" json:"banned"`
	Captcha   *string `form:"captcha" json:"captcha"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
}

type DeleteSelfForm struct {
	Password string `form:"password" json:"password" binding:"required"`
}

type DeleteUserForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}

type DeleteFileForm struct {
	Force bool     `form:"file_ids" json:"force"`
	Files []string `form:"file_ids" json:"file_ids"`
}

type JoinTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Captcha string `form:"captcha" json:"captcha" binding:"required"`
}

type KickMemberForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}

type CreateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
	Desc     string `form:"desc" json:"desc"`
	Country  string `form:"country" json:"country"`
	Hidden   bool   `form:"hidden" json:"hidden"`
	Verified bool   `form:"verified" json:"verified"`
	Banned   bool   `form:"banned" json:"banned"`
}

type CreateTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Captcha string `form:"captcha" json:"captcha"`
}

type CreateContestForm struct {
	Name     string        `form:"name" json:"name" binding:"required"`
	Desc     string        `form:"desc" json:"desc"`
	Start    time.Time     `form:"start" json:"start" binding:"required"`
	Duration time.Duration `form:"duration" json:"duration" binding:"required"`
	Size     int           `form:"size" json:"size" binding:"required"`
	Captcha  string        `form:"captcha" json:"captcha"`
	Hidden   bool          `form:"hidden" json:"hidden"`
}
