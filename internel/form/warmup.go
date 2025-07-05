package form

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"slices"
)

var allowedPullPolicy = []corev1.PullPolicy{corev1.PullAlways, corev1.PullNever, corev1.PullIfNotPresent}

type WarmUpImageForm struct {
	Images     []string `form:"images" json:"images" binding:"required"`
	PullPolicy string   `form:"pull_policy" json:"pull_policy" binding:"required"`
}

func (f *WarmUpImageForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	f.PullPolicy = utils.ToTitle(f.PullPolicy)
	if !slices.Contains(allowedPullPolicy, corev1.PullPolicy(f.PullPolicy)) {
		f.PullPolicy = string(corev1.PullNever)
	}
	return true, i18n.Success
}

type GetContestVictimsForm struct {
	Limit       int    `form:"limit" json:"limit"`
	Offset      int    `form:"offset" json:"offset"`
	ChallengeID string `form:"challenge_id" json:"challenge_id" binding:"omitempty,uuid"`
	TeamID      uint   `form:"team_id" json:"team_id"`
	UserID      uint   `form:"user_id" json:"user_id"`
}

func (f *GetContestVictimsForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	if f.Limit > 100 {
		f.Limit = 100
	}
	if _, exists := ctx.GetQuery("limit"); !exists {
		f.Limit = 10
	}
	if _, exists := ctx.GetQuery("offset"); !exists {
		f.Offset = 0
	}
	return true, i18n.Success
}

type StopContestVictimsForm struct {
	Victims []uint `form:"victims" json:"victims" binding:"required"`
}

func (f *StopContestVictimsForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

type StartContestVictimsForm struct {
	Challenges []string `form:"challenges" json:"challenges" binding:"required"`
	Teams      []uint   `form:"teams" json:"teams" binding:"required"`
}

func (f *StartContestVictimsForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}
