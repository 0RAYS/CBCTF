package form

import (
	"CBCTF/internal/utils"
	"time"
)

// CreateContestForm for create contest
type CreateContestForm struct {
	Name      string               `form:"name" json:"name" binding:"required"`
	Desc      string               `form:"desc" json:"desc"`
	Captcha   string               `form:"captcha" json:"captcha"`
	Prefix    string               `form:"prefix" json:"prefix"`
	Blood     bool                 `form:"blood" json:"blood"`
	Size      int                  `form:"size" json:"size" binding:"required"`
	Start     time.Time            `form:"start" json:"start" binding:"required"`
	Duration  utils.SecondDuration `form:"duration" json:"duration" binding:"required"`
	Rules     utils.Strings        `form:"rules" json:"rules"`
	Prizes    utils.Prizes         `form:"prizes" json:"prizes"`
	Timelines utils.Timelines      `form:"timelines" json:"timelines"`
	Hidden    bool                 `form:"hidden" json:"hidden"`
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name      *string               `form:"name" json:"name"`
	Desc      *string               `form:"desc" json:"desc"`
	Captcha   *string               `form:"captcha" json:"captcha"`
	Blood     *bool                 `form:"blood" json:"blood"`
	Prefix    *string               `form:"prefix" json:"prefix"`
	Size      *int                  `form:"start" json:"size"`
	Start     *time.Time            `form:"start" json:"start"`
	Duration  *utils.SecondDuration `form:"duration" json:"duration"`
	Rules     *utils.Strings        `form:"rules" json:"rules"`
	Prizes    *utils.Prizes         `form:"prizes" json:"prizes"`
	Timelines *utils.Timelines      `form:"timelines" json:"timelines"`
	Hidden    *bool                 `form:"hidden" json:"hidden"`
}
