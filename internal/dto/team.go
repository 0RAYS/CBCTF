package dto

type ListTeamForm struct {
	ListModelsForm
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
}

type CreateTeamForm struct {
	Name        string `form:"name" json:"name" binding:"required"`
	Description string `form:"description" json:"description"`
	Captcha     string `form:"captcha" json:"captcha"`
}

type UpdateTeamForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
	CaptainID   *uint   `form:"captain_id" json:"captain_id"`
}

type JoinTeamForm struct {
	Name    string `form:"name" json:"name" binding:"required"`
	Captcha string `form:"captcha" json:"captcha" binding:"required"`
}

type KickMemberForm struct {
	UserID uint `form:"user_id" json:"user_id" binding:"required"`
}

type AdminUpdateTeamForm struct {
	Name        *string `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string `form:"description" json:"description"`
	Hidden      *bool   `form:"hidden" json:"hidden"`
	Banned      *bool   `form:"banned" json:"banned"`
	Captcha     *string `form:"captcha" json:"captcha" binding:"omitempty,min=1"`
	CaptainID   *uint   `form:"captain_id" json:"captain_id"`
}
