package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"CBCTF/internal/utils"
	"slices"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
)

var allowedPullPolicy = []corev1.PullPolicy{corev1.PullAlways, corev1.PullNever, corev1.PullIfNotPresent}

type WarmUpImageForm struct {
	Images     []string `form:"images" json:"images" binding:"required"`
	PullPolicy string   `form:"pull_policy" json:"pull_policy" binding:"required"`
}

func (f *WarmUpImageForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	f.PullPolicy = utils.ToTitle(f.PullPolicy)
	if !slices.Contains(allowedPullPolicy, corev1.PullPolicy(f.PullPolicy)) {
		f.PullPolicy = string(corev1.PullNever)
	}
	return model.SuccessRetVal()
}

type GetContestVictimsForm struct {
	Limit       int    `form:"limit" json:"limit"`
	Offset      int    `form:"offset" json:"offset"`
	ChallengeID string `form:"challenge_id" json:"challenge_id" binding:"omitempty,uuid"`
	TeamID      uint   `form:"team_id" json:"team_id"`
	UserID      uint   `form:"user_id" json:"user_id"`
}

func (f *GetContestVictimsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	if f.Limit > 100 || f.Limit < 0 {
		f.Limit = 15
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
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
	Challenges []string `form:"challenges" json:"challenges" binding:"required"`
	Teams      []uint   `form:"teams" json:"teams" binding:"required"`
}

func (f *StartContestVictimsForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
