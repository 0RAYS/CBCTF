package model

import (
	"CBCTF/internal/i18n"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

var (
	ContestIsComing  = i18n.Model.Contest.IsComing
	ContestIsRunning = i18n.Model.Contest.IsRunning
	ContestIsOver    = i18n.Model.Contest.IsOver
)

// Contest 赛事
// HasMany Team
// ManyToMany User
// HasMany Notice
// HasMany ContestChallenge
// HasMany ContestFlag
// HasMany Submission
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
	Victims           int64              `gorm:"default:1" json:"victims"`
	Rules             StringList         `gorm:"type:jsonb" json:"rules"`
	Prizes            Prizes             `gorm:"type:jsonb" json:"prizes"`
	Timelines         Timelines          `gorm:"type:jsonb" json:"timelines"`
	BaseModel
}

func (c Contest) IsOver() bool {
	return time.Now().After(c.Start.Add(c.Duration))
}

func (c Contest) IsComing() bool {
	return time.Now().Before(c.Start)
}

func (c Contest) IsRunning() bool {
	return !(c.IsOver() || c.IsComing() || c.Hidden)
}

func (c Contest) Status() string {
	if c.IsOver() {
		return ContestIsOver
	}
	if c.IsComing() {
		return ContestIsComing
	}
	return ContestIsRunning
}

type Prize struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
}

type Prizes []Prize

func (p Prizes) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Prizes) Scan(value any) error {
	bytes, err := scanBytes(value)
	if err != nil {
		return fmt.Errorf("failed to scan Prizes value")
	}
	return json.Unmarshal(bytes, p)
}

type Timeline struct {
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

type Timelines []Timeline

func (t Timelines) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Timelines) Scan(value any) error {
	bytes, err := scanBytes(value)
	if err != nil {
		return fmt.Errorf("failed to scan Timelines value")
	}
	return json.Unmarshal(bytes, t)
}
