package dto

import (
	"CBCTF/internal/model"
)

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
