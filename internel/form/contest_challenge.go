package form

import (
	"CBCTF/internel/i18n"
	"CBCTF/internel/model"
	"CBCTF/internel/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

// CreateContestChallengeForm add challenge to contest
type CreateContestChallengeForm struct {
	ChallengeRandIDL []string `form:"challenge_id" json:"challenge_id" binding:"required"`
}

func (f *CreateContestChallengeForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	return true, i18n.Success
}

type UpdateContestChallengeForm struct {
	Name    *string           `form:"name" json:"name"`
	Desc    *string           `form:"desc" json:"desc"`
	Hidden  *bool             `form:"hidden" json:"hidden"`
	Attempt *int64            `form:"attempt" json:"attempt"`
	Hints   *model.StringList `form:"hints" json:"hints"`
	Tags    *model.StringList `form:"tags" json:"tags"`
}

func (f *UpdateContestChallengeForm) Bind(ctx *gin.Context) (bool, string) {
	if err := ctx.ShouldBind(f); err != nil {
		return false, i18n.BadRequest
	}
	if f.Name != nil {
		f.Name = utils.Ptr(strings.TrimSpace(*f.Name))
		if *f.Name == "" {
			return false, i18n.BadRequest
		}
	}
	return true, i18n.Success
}
