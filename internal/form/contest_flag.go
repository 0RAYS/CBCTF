package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"
	"slices"

	"github.com/gin-gonic/gin"
)

var allowedScoreType = []uint{model.StaticScore, model.LinearScore, model.LogarithmicScore}

type UpdateContestFlagForm struct {
	Value     *string  `form:"value" json:"value"`
	Score     *float64 `form:"score" json:"score"`
	Decay     *float64 `form:"decay" json:"decay"`
	MinScore  *float64 `form:"min_score" json:"min_score"`
	ScoreType *uint    `form:"score_type" json:"score_type"`
}

func (f *UpdateContestFlagForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %v", err)
		return false, i18n.BadRequest
	}
	if f.ScoreType != nil {
		if !slices.Contains(allowedScoreType, *f.ScoreType) {
			return false, i18n.InvalidScoreType
		}
	}
	return true, i18n.Success
}
