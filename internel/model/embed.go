package model

import (
	"CBCTF/internel/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type AvatarURL string

func (a AvatarURL) Value() (driver.Value, error) {
	return strings.TrimPrefix(string(a), strings.Trim(config.Env.Backend, "/")), nil
}

func (a *AvatarURL) Scan(value any) error {
	path, ok := value.(string)
	if !ok {
		return fmt.Errorf("failed to scan AvatarURL: %v", value)
	}
	*a = AvatarURL(strings.Trim(config.Env.Backend, "/") + path)
	return nil
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *StringList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringList value")
	}
	return json.Unmarshal(bytes, s)
}

type UintList []uint

func (u UintList) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *UintList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UintList value")
	}
	return json.Unmarshal(bytes, u)
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

type Reference struct {
	UserID    uint `json:"user_id"`
	TeamID    uint `json:"team_id"`
	ContestID uint `json:"contest_id"`
	UsageID   uint `json:"usage_id"`
}

func (r Reference) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *Reference) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Reference value")
	}
	return json.Unmarshal(bytes, r)
}
