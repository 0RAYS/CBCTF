package dto

import (
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type PullImageForm struct {
	Images     []string `form:"images" json:"images" binding:"required,dive,min=1"`
	PullPolicy string   `form:"pull_policy" json:"pull_policy" binding:"required,oneof=Always Never IfNotPresent"`
}

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
	Teams      []uint   `form:"teams" json:"teams" binding:"required,dive,gt=0"`
}
