package model

import (
	"CBCTF/internal/config"
	"CBCTF/internal/form"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Rules []string

func (r Rules) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Rules) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Rules value")
	}
	return json.Unmarshal(bytes, r)
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

type Contest struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"index:idx_name_deleted,unique;not null" json:"name"`
	Desc      string         `json:"desc"`
	Captcha   string         `json:"-"`
	Avatar    string         `json:"avatar"`
	Prefix    string         `json:"prefix" gorm:"default:'CBCTF'"`
	Size      int            `json:"size"`
	Start     time.Time      `json:"start"`
	Duration  time.Duration  `json:"-"`
	Blood     bool           `json:"blood" gorm:"default:true"`
	Hidden    bool           `gorm:"default:true" json:"hidden"`
	Rules     Rules          `json:"rules" gorm:"type:json"`
	Prizes    Prizes         `json:"prizes" gorm:"type:json"`
	Timelines Timelines      `json:"timelines" gorm:"type:json"`
	Teams     []*Team        `json:"-"`
	Users     []*User        `gorm:"many2many:user_contests;" json:"-"`
	Notices   []*Notice      `json:"-"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index;index:idx_name_deleted,unique;" json:"-"`
}

func (c *Contest) MarshalJSON() ([]byte, error) {
	type Tmp Contest // 定义一个别名以避免递归调用
	return json.Marshal(&struct {
		*Tmp
		Users    int    `json:"users"`
		Teams    int    `json:"teams"`
		Notices  int    `json:"notices"`
		Avatar   string `json:"avatar"`
		Duration int64  `json:"duration"`
	}{
		Tmp:      (*Tmp)(c),
		Users:    len(c.Users),
		Teams:    len(c.Teams),
		Notices:  len(c.Notices),
		Avatar:   fmt.Sprintf("%s/%s", config.Env.Backend, strings.TrimPrefix(c.Avatar, "/")),
		Duration: int64(c.Duration.Seconds()),
	})
}

func (c *Contest) IsOver() bool {
	return time.Now().After(c.Start.Add(c.Duration))
}

func (c *Contest) IsNotStart() bool {
	return time.Now().Before(c.Start)
}

func (c *Contest) IsRunning() bool {
	return (c.IsOver() || c.IsNotStart()) != true
}

func (c *Contest) Status() string {
	if c.IsOver() {
		return "ContestIsOver"
	}
	if c.IsNotStart() {
		return "ContestNotRunning"
	}
	return "ContestIsRunning"
}

func InitContest(form form.CreateContestForm) Contest {
	return Contest{
		Name:     form.Name,
		Desc:     form.Desc,
		Captcha:  form.Captcha,
		Avatar:   "",
		Blood:    form.Blood,
		Size:     form.Size,
		Start:    form.Start,
		Hidden:   form.Hidden,
		Duration: time.Duration(form.Duration),
	}
}
