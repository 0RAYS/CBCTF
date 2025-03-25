package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

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

type Uints []uint

func (u Uints) Value() (driver.Value, error) {
	return json.Marshal(u)
}

func (u *Uints) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Uints value")
	}
	return json.Unmarshal(bytes, u)
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

type Docker struct {
	Flags Strings `json:"flags"`
	Image string  `json:"image"`
	Ports Uints   `json:"ports"`
}

func (d Docker) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Docker) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Docker value")
	}
	return json.Unmarshal(bytes, d)
}

type Dockers []Docker

func (d Dockers) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *Dockers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Dockers value")
	}
	return json.Unmarshal(bytes, d)
}

type Expose struct {
	IP   string `json:"ip"`
	Port int32  `json:"port"`
}

type Exposes []Expose

func (e Exposes) Value() (driver.Value, error) {
	return json.Marshal(e)
}

func (e *Exposes) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Exposes value")
	}
	return json.Unmarshal(bytes, e)
}
