package form

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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

type Strings []string

func (s Strings) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *Strings) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Strings value")
	}
	return json.Unmarshal(bytes, s)
}

type Prize struct {
	Amount string `json:"amount"`
	Desc   string `json:"desc"`
}

type Prizes []Prize

func (p Prizes) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Prizes) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Prizes value")
	}
	return json.Unmarshal(bytes, p)
}

type Timeline struct {
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
	Desc  string    `json:"desc"`
}

type Timelines []Timeline

func (t Timelines) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Timelines) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Timelines value")
	}
	return json.Unmarshal(bytes, t)
}

// CreateContestForm for create contest
type CreateContestForm struct {
	Name      string         `form:"name" json:"name" binding:"required"`
	Desc      string         `form:"desc" json:"desc"`
	Captcha   string         `form:"captcha" json:"captcha"`
	Prefix    string         `form:"prefix" json:"prefix"`
	Blood     bool           `form:"blood" json:"blood"`
	Size      int            `form:"size" json:"size" binding:"required"`
	Start     time.Time      `form:"start" json:"start" binding:"required"`
	Duration  SecondDuration `form:"duration" json:"duration" binding:"required"`
	Rules     Strings        `form:"rules" json:"rules"`
	Prizes    Prizes         `form:"prizes" json:"prizes"`
	Timelines Timelines      `form:"timelines" json:"timelines"`
	Hidden    bool           `form:"hidden" json:"hidden"`
}

// UpdateContestForm for admin update contest info
type UpdateContestForm struct {
	Name      *string         `form:"name" json:"name"`
	Desc      *string         `form:"desc" json:"desc"`
	Captcha   *string         `form:"captcha" json:"captcha"`
	Blood     *bool           `form:"blood" json:"blood"`
	Prefix    *string         `form:"prefix" json:"prefix"`
	Size      *int            `form:"start" json:"size"`
	Start     *time.Time      `form:"start" json:"start"`
	Duration  *SecondDuration `form:"duration" json:"duration"`
	Rules     *Strings        `form:"rules" json:"rules"`
	Prizes    *Prizes         `form:"prizes" json:"prizes"`
	Timelines *Timelines      `form:"timelines" json:"timelines"`
	Hidden    *bool           `form:"hidden" json:"hidden"`
}
