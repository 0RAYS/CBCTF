package form

import (
	"CBCTF/internal/i18n"
	"CBCTF/internal/log"
	"CBCTF/internal/model"

	"github.com/gin-gonic/gin"
)

// CreateContestChallengeForm add challenge to contest
type CreateContestChallengeForm struct {
	ChallengeRandIDL []string `form:"challenge_id" json:"challenge_id" binding:"required"`
}

func (f *CreateContestChallengeForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}

type UpdateContestChallengeForm struct {
	Name    *string           `form:"name" json:"name"`
	Desc    *string           `form:"desc" json:"desc"`
	Hidden  *bool             `form:"hidden" json:"hidden"`
	Attempt *int64            `form:"attempt" json:"attempt"`
	Hints   *model.StringList `form:"hints" json:"hints"`
	Tags    *model.StringList `form:"tags" json:"tags"`
}

func (f *UpdateContestChallengeForm) Bind(ctx *gin.Context) model.RetVal {
	if err := ctx.ShouldBind(f); err != nil {
		log.Logger.Debugf("Failed to bind form: %s", err)
		return model.RetVal{Msg: i18n.Request.BadRequest, Attr: map[string]any{"Error": err.Error()}}
	}
	return model.SuccessRetVal()
}
