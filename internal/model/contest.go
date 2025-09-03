package model

import (
	"CBCTF/internal/i18n"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

var (
	ContestIsComing  = i18n.ContestIsComing
	ContestIsRunning = i18n.ContestIsRunning
	ContestIsOver    = i18n.ContestIsOver
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
	Users             []*User            `gorm:"many2many:user_contests;" json:"-"`
	Notices           []Notice           `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContestChallenges []ContestChallenge `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	ContestFlags      []ContestFlag      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Submissions       []Submission       `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Name              string             `gorm:"type:varchar(255);uniqueIndex;not null" json:"name"`
	Desc              string             `json:"desc"`
	Captcha           string             `json:"captcha"`
	Avatar            AvatarURL          `json:"avatar"`
	Prefix            string             `gorm:"default:'flag'" json:"prefix"`
	Size              int                `gorm:"default:4" json:"size"`
	Start             time.Time          `gorm:"not null" json:"start"`
	Duration          time.Duration      `json:"duration"`
	Blood             bool               `gorm:"default:true" json:"blood"`
	Hidden            bool               `gorm:"default:true" json:"hidden"`
	Rules             StringList         `gorm:"type:json" json:"rules"`
	Prizes            Prizes             `gorm:"type:json" json:"prizes"`
	Timelines         Timelines          `gorm:"type:json" json:"timelines"`
	BasicModel
}

func (c Contest) GetModelName() string {
	return "Contest"
}

func (c Contest) GetBasicModel() BasicModel {
	return c.BasicModel
}

func (c Contest) CreateErrorString() string {
	return i18n.CreateContestError
}

func (c Contest) DeleteErrorString() string {
	return i18n.DeleteContestError
}

func (c Contest) GetErrorString() string {
	return i18n.GetContestError
}

func (c Contest) NotFoundErrorString() string {
	return i18n.ContestNotFound
}

func (c Contest) UpdateErrorString() string {
	return i18n.UpdateContestError
}

func (c Contest) GetUniqueKey() []string {
	return []string{"id", "name"}
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

type Prizes []struct {
	Amount string `json:"amount"`
	Desc   string `json:"desc"`
}

func (p Prizes) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *Prizes) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Prizes value")
	}
	return json.Unmarshal(bytes, p)
}

type Timelines []struct {
	Date  time.Time `json:"date"`
	Title string    `json:"title"`
	Desc  string    `json:"desc"`
}

func (t Timelines) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Timelines) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Timelines value")
	}
	return json.Unmarshal(bytes, t)
}
