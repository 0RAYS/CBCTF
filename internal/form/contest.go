package form

import "time"

// CreateContestForm for create contest
type CreateContestForm struct {
	Name     string        `form:"name" json:"name" binding:"required"`
	Desc     string        `form:"desc" json:"desc"`
	Captcha  string        `form:"captcha" json:"captcha"`
	Prefix   string        `form:"prefix" json:"prefix"`
	Blood    bool          `form:"blood" json:"blood"`
	Size     int           `form:"size" json:"size" binding:"required"`
	Start    time.Time     `form:"start" json:"start" binding:"required"`
	Duration time.Duration `form:"duration" json:"duration" binding:"required"`
	Hidden   bool          `form:"hidden" json:"hidden"`
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name     *string        `form:"name" json:"name"`
	Desc     *string        `form:"desc" json:"desc"`
	Captcha  *string        `form:"captcha" json:"captcha"`
	Blood    *bool          `form:"blood" json:"blood"`
	Prefix   *string        `form:"prefix" json:"prefix"`
	Avatar   *string        `form:"avatar" json:"avatar"`
	Size     *int           `form:"start" json:"size"`
	Start    *time.Time     `form:"start" json:"start"`
	Duration *time.Duration `form:"duration" json:"duration"`
	Hidden   *bool          `form:"hidden" json:"hidden"`
}
