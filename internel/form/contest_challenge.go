package form

import "CBCTF/internel/model"

// CreateContestChallengeForm add challenge to contest
type CreateContestChallengeForm struct {
	ChallengeRandIDL []string `form:"challenge_id" json:"challenge_id" binding:"required"`
}

type UpdateContestChallengeForm struct {
	Name    *string           `form:"name" json:"name"`
	Desc    *string           `form:"desc" json:"desc"`
	Hidden  *bool             `form:"hidden" json:"hidden"`
	Attempt *int64            `form:"attempt" json:"attempt"`
	Hints   *model.StringList `form:"hints" json:"hints"`
	Tags    *model.StringList `form:"tags" json:"tags"`
}
