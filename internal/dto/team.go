package dto

import (
	"strings"

	"github.com/gin-gonic/gin"

	"CBCTF/internal/i18n"
	"CBCTF/internal/model"
)

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

// UpdateTeamScoreForm 管理员对队伍分数的人工变动。Score 是带符号的变动分数，Reason 是变动原因。
type UpdateTeamScoreForm struct {
	Score  float64 `form:"score" json:"score"`
	Reason string  `form:"reason" json:"reason"`
}

func (f *UpdateTeamScoreForm) Validate(_ *gin.Context) model.RetVal {
	f.Reason = strings.TrimSpace(f.Reason)
	if f.Score == 0 {
		return model.FailRetVal(i18n.Model.Team.ScoreInvalid, "Team", "score", nil)
	}
	if f.Reason == "" {
		return model.FailRetVal(i18n.Model.Team.ScoreReasonRequired, "Team", "reason", nil)
	}
	return model.SuccessRetVal()
}
