package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type WarmUpImageForm struct {
	Images     []string `form:"images" json:"images" binding:"required"`
	PullPolicy string   `form:"pull_policy" json:"pull_policy" binding:"required,oneof=Always Never IfNotPresent"`
}

func (f *WarmUpImageForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type GetContestVictimsForm struct {
	Offset      int    `form:"offset" json:"offset" binding:"gte=0"`
	Limit       int    `form:"limit" json:"limit" binding:"gte=0,lte=100"`
	ChallengeID string `form:"challenge_id" json:"challenge_id" binding:"omitempty,uuid"`
	TeamID      uint   `form:"team_id" json:"team_id"`
	UserID      uint   `form:"user_id" json:"user_id"`
}

func (f *GetContestVictimsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type StopContestVictimsForm struct {
	Victims []uint `form:"victims" json:"victims" binding:"required"`
}

func (f *StopContestVictimsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type StartContestVictimsForm struct {
	Challenges []string `form:"challenges" json:"challenges" binding:"required,dive,uuid"`
	Teams      []uint   `form:"teams" json:"teams" binding:"required"`
}

func (f *StartContestVictimsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
