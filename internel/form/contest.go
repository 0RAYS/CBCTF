package form

import (
	"CBCTF/internel/model"
	"time"
)

// CreateContestForm for create contest
type CreateContestForm struct {
	Name      string           `form:"name" json:"name" binding:"required"`
	Desc      string           `form:"desc" json:"desc"`
	Captcha   string           `form:"captcha" json:"captcha"`
	Prefix    string           `form:"prefix" json:"prefix"`
	Blood     bool             `form:"blood" json:"blood"`
	Size      int              `form:"size" json:"size"`
	Start     time.Time        `form:"start" json:"start"`
	Duration  int64            `form:"duration" json:"duration"`
	Rules     model.StringList `form:"rules" json:"rules"`
	Prizes    model.Prizes     `form:"prizes" json:"prizes"`
	Timelines model.Timelines  `form:"timelines" json:"timelines"`
	Hidden    bool             `form:"hidden" json:"hidden"`
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name      *string           `form:"name" json:"name"`
	Desc      *string           `form:"desc" json:"desc"`
	Captcha   *string           `form:"captcha" json:"captcha"`
	Blood     *bool             `form:"blood" json:"blood"`
	Prefix    *string           `form:"prefix" json:"prefix"`
	Size      *int              `form:"start" json:"size"`
	Start     *time.Time        `form:"start" json:"start"`
	Duration  *int64            `form:"duration" json:"duration"`
	Rules     *model.StringList `form:"rules" json:"rules"`
	Prizes    *model.Prizes     `form:"prizes" json:"prizes"`
	Timelines *model.Timelines  `form:"timelines" json:"timelines"`
	Hidden    *bool             `form:"hidden" json:"hidden"`
}
