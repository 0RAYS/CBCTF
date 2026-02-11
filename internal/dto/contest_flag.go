package dto

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

type UpdateContestFlagForm struct {
	Value     *string  `form:"value" json:"value" binding:"omitempty,min=1"`
	Score     *float64 `form:"score" json:"score" binding:"omitempty,gte=0"`
	Decay     *float64 `form:"decay" json:"decay" binding:"omitempty,gte=0"`
	MinScore  *float64 `form:"min_score" json:"min_score" binding:"omitempty,gte=0"`
	ScoreType *uint    `form:"score_type" json:"score_type" binding:"omitempty,oneof=0 1 2"`
}

func (f *UpdateContestFlagForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
