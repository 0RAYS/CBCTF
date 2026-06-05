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
	Duration   int64    `form:"duration" json:"duration" binding:"omitempty,gte=1"`
}

type GetVictimPodLogsForm struct {
	PodName   string `form:"pod_name" json:"pod_name" binding:"required"`
	Container string `form:"container" json:"container" binding:"required"`
	Lines     int64  `form:"lines" json:"lines" binding:"omitempty,gt=0"`
}

func (f *GetVictimPodLogsForm) Validate(_ *gin.Context) model.RetVal {
	if f.Lines <= 0 {
		f.Lines = 1000
	}
	return model.SuccessRetVal()
}
