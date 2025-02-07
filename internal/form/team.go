package form

// CreateTeamForm for create team
type CreateTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Desc    string `form:"desc" json:"desc"`
	Captcha string `form:"captcha" json:"captcha"`
}

// UpdateTeamForm for user update team info
type UpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	Captcha   *string `form:"captcha" json:"captcha"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
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

// AdminUpdateTeamForm for admin update team info
type AdminUpdateTeamForm struct {
	Name      *string `form:"name" json:"name"`
	Desc      *string `form:"desc" json:"desc"`
	Hidden    *bool   `form:"hidden" json:"hidden"`
	Banned    *bool   `form:"banned" json:"banned"`
	Captcha   *string `form:"captcha" json:"captcha"`
	CaptainID *uint   `form:"captain_id" json:"captain_id"`
}
