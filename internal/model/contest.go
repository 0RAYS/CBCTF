package model

import (
	"time"

	"CBCTF/internal/i18n"
)

var (
	ContestIsComing  = i18n.Model.Contest.IsComing
	ContestIsRunning = i18n.Model.Contest.IsRunning
	ContestIsPaused  = i18n.Model.Contest.IsPaused
	ContestIsOver    = i18n.Model.Contest.IsOver
)

type Contest struct {
	Teams             []Team             `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Users             []User             `gorm:"many2many:user_contests;" json:"-"`
	Notices           []Notice           `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContestChallenges []ContestChallenge `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContestFlags      []ContestFlag      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions       []Submission       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name              string             `gorm:"type:varchar(255);uniqueIndex:idx_contests_name_active,where:deleted_at IS NULL;not null" json:"name"`
	Description       string             `json:"description"`
	Captcha           string             `json:"captcha"`
	Picture           FileURL            `json:"picture"`
	Prefix            string             `gorm:"default:'flag'" json:"prefix"`
	Size              int                `gorm:"default:4" json:"size"`
	Start             time.Time          `gorm:"not null" json:"start"`
	Duration          time.Duration      `json:"duration"`
	Blood             bool               `gorm:"default:true" json:"blood"`
	Hidden            bool               `gorm:"default:true;index" json:"hidden"`
	Paused            bool               `gorm:"default:false" json:"paused"`
	PausedAt          time.Time          `gorm:"default:null" json:"-"`
	Victims           int64              `gorm:"default:1" json:"victims"`
	Rules             []string           `gorm:"serializer:json;type:jsonb;default:'[]'" json:"rules"`
	Prizes            []Prize            `gorm:"serializer:json;type:jsonb;default:'[]'" json:"prizes"`
	Timelines         []Timeline         `gorm:"serializer:json;type:jsonb;default:'[]'" json:"timelines"`
	BaseModel
}

// End returns the moment the contest stops scoring.
// Time spent paused is not part of the contest clock.
func (c Contest) End() time.Time {
	end := c.Start.Add(c.Duration)
	if c.Paused && !c.PausedAt.IsZero() {
		end = end.Add(time.Since(c.PausedAt))
	}
	return end
}

// EndedWithin reports whether the contest ended within the given window
// before now. A contest that has not ended yet is always inside the window.
func (c Contest) EndedWithin(now time.Time, window time.Duration) bool {
	return now.Sub(c.End()) <= window
}

func (c Contest) IsOver() bool {
	return time.Now().After(c.End())
}

func (c Contest) IsComing() bool {
	return time.Now().Before(c.Start)
}

func (c Contest) IsRunning() bool {
	return !(c.IsOver() || c.IsComing() || c.Hidden || c.Paused)
}

func (c Contest) Status() string {
	if c.IsOver() {
		return ContestIsOver
	}
	if c.IsComing() {
		return ContestIsComing
	}
	if c.Paused {
		return ContestIsPaused
	}
	return ContestIsRunning
}

type Prize struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

type Timeline struct {
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}
