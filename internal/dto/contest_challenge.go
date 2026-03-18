package dto

import (
	"CBCTF/internal/model"
	"CBCTF/internal/utils"

	"github.com/gin-gonic/gin"
)

type GetContestChallengesForm struct {
	ListModelsForm
	Category string `form:"category" json:"category"`
}

func (f *GetContestChallengesForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

type GetAllContestChallengesForm struct {
	ListModelsForm
	Type     model.ChallengeType `form:"type" json:"type" binding:"omitempty,oneof=static question dynamic pods"`
	Category string              `form:"category" json:"category"`
}

func (f *GetAllContestChallengesForm) Validate(ctx *gin.Context) model.RetVal {
	if _, ok := ctx.GetQuery("limit"); !ok {
		f.Limit = 10
	}
	if _, ok := ctx.GetQuery("offset"); !ok {
		f.Offset = 0
	}
	f.Category = utils.ToTitle(f.Category)
	return model.SuccessRetVal()
}

// CreateContestChallengeForm add challenge to contest
type CreateContestChallengeForm struct {
	ChallengeIDs []string `form:"challenge_ids" json:"challenge_ids" binding:"required,dive,uuid"`
}

type UpdateContestChallengeForm struct {
	Name        *string           `form:"name" json:"name" binding:"omitempty,min=1"`
	Description *string           `form:"description" json:"description"`
	Hidden      *bool             `form:"hidden" json:"hidden"`
	Attempt     *int64            `form:"attempt" json:"attempt" binding:"omitempty,gte=0"`
	Hints       *model.StringList `form:"hints" json:"hints"`
	Tags        *model.StringList `form:"tags" json:"tags"`
}
