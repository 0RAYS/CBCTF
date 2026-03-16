package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type GetVictimsForm struct {
	ListModelsForm
	ChallengeID string `form:"challenge_id" json:"challenge_id" binding:"omitempty,uuid"`
	TeamID      uint   `form:"team_id" json:"team_id"`
	UserID      uint   `form:"user_id" json:"user_id"`
	Deleted     bool   `form:"deleted" json:"deleted"`
}

func (f *GetVictimsForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	return model.SuccessRetVal()
}

type StopVictimsForm struct {
	Victims []uint `form:"victims" json:"victims" binding:"required,dive,gt=0"`
}

type StartVictimsForm struct {
	Challenges []string `form:"challenges" json:"challenges" binding:"required,dive,uuid"`
	TeamRatio  float64  `form:"team_ratio" json:"team_ratio" binding:"required,gt=0,lt=1"`
}
