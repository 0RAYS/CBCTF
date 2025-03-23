package utils

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
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

func (p NetworkPolicy) Value() (driver.Value, error) {
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
	return json.Marshal(p)
}

func (p *NetworkPolicy) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	return json.Unmarshal(bytes, p)
}

type NetworkPolicies []NetworkPolicy

func (p NetworkPolicies) Value() (driver.Value, error) {
	for _, policy := range p {
		for i, ipBlock := range policy.From {
			if !isValidIPBlock(ipBlock) {
				policy.From = append(policy.From[:i], policy.From[i+1:]...)
			}
		}
		for i, ipBlock := range policy.To {
			if !isValidIPBlock(ipBlock) {
				policy.To = append(policy.To[:i], policy.To[i+1:]...)
			}
		}
	}
	return json.Marshal(p)
}

func (p *NetworkPolicies) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicies value")
	}
	return json.Unmarshal(bytes, p)
}

type Docker struct {
	Image string  `json:"image"`
	Ports []int32 `json:"ports"`
}

type Dockers []Docker

func (d Dockers) Value() (driver.Value, error) {
	for _, docker := range d {
		for i, port := range docker.Ports {
			if port < 0 || port > 65535 {
				docker.Ports = append(docker.Ports[:i], docker.Ports[i+1:]...)
			}
		}
	}
	return json.Marshal(d)
}

func (d *Dockers) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Dockers value")
	}
	return json.Unmarshal(bytes, d)
}
