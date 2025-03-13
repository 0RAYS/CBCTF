package utils

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
