package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
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

type IPBlock struct {
	CIDR   string   `json:"cidr"`
	Except []string `json:"except"`
}

func isValidIPBlock(ipBlock IPBlock) bool {
	_, ipNet, err := net.ParseCIDR(ipBlock.CIDR)
	if err != nil {
		return false
	}
	for _, ex := range ipBlock.Except {
		_, exNet, err := net.ParseCIDR(ex)
		if err != nil {
			return false
		}
		if !ipNet.Contains(exNet.IP) {
			return false
		}
	}
	return true
}

type NetworkPolicy struct {
	From []IPBlock `json:"from"`
	To   []IPBlock `json:"to"`
}

type NetworkPolicies []NetworkPolicy

func (n NetworkPolicies) Value() (driver.Value, error) {
	for _, p := range n {
		for i, ipBlock := range p.From {
			if !isValidIPBlock(ipBlock) {
				p.From = append(p.From[:i], p.From[i+1:]...)
			}
		}
		for i, ipBlock := range p.To {
			if !isValidIPBlock(ipBlock) {
				p.To = append(p.To[:i], p.To[i+1:]...)
			}
		}
	}
	return json.Marshal(n)
}

func (n *NetworkPolicies) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	return json.Unmarshal(bytes, n)
}
