package form

import (
	"CBCTF/internel/model"
)

// CreateUsageForm add challenge to contest
type CreateUsageForm struct {
	ChallengeIDL []string `form:"challenge_id" json:"challenge_id" binding:"required"`
}

// UpdateUsageForm for admin update usage info
type UpdateUsageForm struct {
	Name    *string        `form:"name" json:"name"`
	Desc    *string        `form:"desc" json:"desc"`
	Attempt *int64         `form:"desc" json:"attempt"`
	Hidden  *bool          `form:"hidden" json:"hidden"`
	Hints   *model.Strings `form:"hints" json:"hints"`
	Tags    *model.Strings `form:"tags" json:"tags"`
	Docker  *model.Docker  `form:"docker" json:"docker"`
	Dockers *model.Dockers `form:"dockers" json:"dockers"`
}
