package constants

import "time"

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

// CreateAdminForm for create admin
type CreateAdminForm RegisterForm

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

// CreateTeamForm for create team
type CreateTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Desc    string `form:"desc" json:"desc"`
	Captcha string `form:"captcha" json:"captcha"`
}

// CreateContestForm for create contest
type CreateContestForm struct {
	Name     string        `form:"name" json:"name" binding:"required"`
	Desc     string        `form:"desc" json:"desc"`
	Captcha  string        `form:"captcha" json:"captcha"`
	Prefix   string        `form:"prefix" json:"prefix"`
	Size     int           `form:"size" json:"size" binding:"required"`
	Start    time.Time     `form:"start" json:"start" binding:"required"`
	Duration time.Duration `form:"duration" json:"duration" binding:"required"`
	Hidden   bool          `form:"hidden" json:"hidden"`
}

// CreateChallengeForm for create challenge
type CreateChallengeForm struct {
	Name           string `form:"name" json:"name" binding:"required"`
	Desc           string `form:"desc" json:"desc"`
	Flag           string `form:"flag" json:"flag"`
	Category       string `form:"category" json:"category"`
	Type           int    `form:"type" json:"type"`
	GeneratorImage string `form:"generator" json:"generator"`
	DockerImage    string `form:"docker" json:"docker"`
	Port           int32  `form:"port" json:"port"`
}

type CreateUsageForm struct {
	ChallengeID []string `form:"challenge_id" json:"challenge_id" binding:"required"`
}

// GetModelsForm for get models list
type GetModelsForm struct {
	Offset int `form:"offset" json:"offset"`
	Limit  int `form:"limit" json:"limit"`
}

type GetChallengesForm struct {
	Offset   int    `form:"offset" json:"offset"`
	Limit    int    `form:"limit" json:"limit"`
	Type     int    `form:"type" json:"type"`
	Category string `form:"category" json:"category"`
}

// ChangePasswordForm for user or admin change password
type ChangePasswordForm struct {
	OldPassword string `form:"oldPassword" json:"oldPassword" binding:"required"`
	NewPassword string `form:"newPassword" json:"newPassword" binding:"required"`
}

// UpdateSelfForm for user update info
type UpdateSelfForm struct {
	Name    *string `form:"name" json:"name"`
	Email   *string `form:"email" json:"email"`
	Desc    *string `form:"desc" json:"desc"`
	Country *string `form:"country" json:"country"`
}

// UpdateAdminForm for admin update info
type UpdateAdminForm struct {
	Name  *string `form:"name" json:"name"`
	Email *string `form:"email" json:"email"`
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

// UpdateTeamForm for user update team info
type UpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	Captcha   *string `form:"captcha" json:"captcha"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
}

// AdminUpdateTeamForm for admin update team info
type AdminUpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	Hidden    *bool   `form:"hidden" json:"hidden"`
	Banned    *bool   `form:"banned" json:"banned"`
	Captcha   *string `form:"captcha" json:"captcha"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name     *string        `form:"name" json:"name"`
	Desc     *string        `form:"desc" json:"desc"`
	Captcha  *string        `form:"captcha" json:"captcha"`
	Prefix   *string        `form:"prefix" json:"prefix"`
	Avatar   *string        `form:"avatar" json:"avatar"`
	Size     *int           `form:"start" json:"size"`
	Start    *time.Time     `form:"start" json:"start"`
	Duration *time.Duration `form:"duration" json:"duration"`
	Hidden   *bool          `form:"hidden" json:"hidden"`
}

// UpdateChallengeForm for admin update challenge info
type UpdateChallengeForm struct {
	Name           *string `form:"name" json:"name"`
	Desc           *string `form:"desc" json:"desc"`
	Flag           *string `form:"flag" json:"flag"`
	Category       *string `form:"category" json:"category"`
	Type           *int    `form:"type" json:"type"`
	GeneratorImage *string `form:"generator" json:"generator"`
	DockerImage    *string `form:"docker" json:"docker"`
	Port           *int32  `form:"port" json:"port"`
}

// UpdateUsageForm for admin update usage info
type UpdateUsageForm struct {
	Hidden  *bool   `form:"hidden" json:"hidden"`
	Score   *int    `form:"score" json:"score"`
	Attempt *int    `form:"attempt" json:"attempt"`
	Hints   *string `form:"hints" json:"hints"`
	Tags    *string `form:"tags" json:"tags"`
}

// DeleteSelfForm for user delete self
type DeleteSelfForm struct {
	Password string `form:"password" json:"password" binding:"required"`
}

// DeleteUserForm for admin delete user
type DeleteUserForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}

// DeleteFileForm for delete files
type DeleteFileForm struct {
	Force bool     `form:"force" json:"force"`
	Files []string `form:"file_ids" json:"file_ids"`
}

type DeleteChallengeForm struct {
	Force bool `form:"force" json:"force"`
}

// JoinTeamForm for user join team
type JoinTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Captcha string `form:"captcha" json:"captcha" binding:"required"`
}

// KickMemberForm for admin or captain kick member
type KickMemberForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}

type GetCategoriesForm struct {
	Type int `form:"type" json:"type"`
}

type DownloadChallengeForm struct {
	File string `form:"file" json:"file" binding:"required"`
}

type SubmitFlagForm struct {
	Flag string `form:"flag" json:"flag" binding:"required"`
}
