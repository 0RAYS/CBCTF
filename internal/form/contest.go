package form

import (
	"strconv"
	"time"
)

type SecondDuration time.Duration

func (d *SecondDuration) UnmarshalJSON(b []byte) error {
	seconds, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*d = SecondDuration(time.Duration(seconds) * time.Second)
	return nil
}

func (d *SecondDuration) UnmarshalText(b []byte) error {
	seconds, err := strconv.ParseInt(string(b), 10, 64)
	if err != nil {
		return err
	}
	*d = SecondDuration(time.Duration(seconds) * time.Second)
	return nil
}

// CreateContestForm for create contest
type CreateContestForm struct {
	Name     string         `form:"name" json:"name" binding:"required"`
	Desc     string         `form:"desc" json:"desc"`
	Captcha  string         `form:"captcha" json:"captcha"`
	Prefix   string         `form:"prefix" json:"prefix"`
	Blood    bool           `form:"blood" json:"blood"`
	Size     int            `form:"size" json:"size" binding:"required"`
	Start    time.Time      `form:"start" json:"start" binding:"required"`
	Duration SecondDuration `form:"duration" json:"duration" binding:"required"`
	Hidden   bool           `form:"hidden" json:"hidden"`
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name     *string         `form:"name" json:"name"`
	Desc     *string         `form:"desc" json:"desc"`
	Captcha  *string         `form:"captcha" json:"captcha"`
	Blood    *bool           `form:"blood" json:"blood"`
	Prefix   *string         `form:"prefix" json:"prefix"`
	Size     *int            `form:"start" json:"size"`
	Start    *time.Time      `form:"start" json:"start"`
	Duration *SecondDuration `form:"duration" json:"duration"`
	Hidden   *bool           `form:"hidden" json:"hidden"`
}
