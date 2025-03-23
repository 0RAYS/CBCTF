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

var DefaultNetworkPolicy = NetworkPolicy{
	From: []IPBlock{},
	To: []IPBlock{
		{
			CIDR: "0.0.0.0/0",
			Except: []string{
				"10.0.0.0/8",
				"172.16.0.0/12",
				"192.168.0.0/16",
				"100.64.0.0/10",
			},
		},
	},
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

type IPPort struct {
	IP   string `json:"ip"`
	Port int32  `json:"port"`
}

type IPPorts []IPPort

func (i IPPorts) Value() (driver.Value, error) {
	for _, ipPort := range i {
		if ipPort.Port < 0 || ipPort.Port > 65535 {
			return nil, fmt.Errorf("invalid port")
		}
	}
	return json.Marshal(i)
}

func (i *IPPorts) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan IPPorts value")
	}
	return json.Unmarshal(bytes, i)
}

type Docker struct {
	Image         string        `json:"image"`
	Flag          string        `json:"flag"`
	Ports         []int32       `json:"ports"`
	NetworkPolicy NetworkPolicy `json:"network_policy"`
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
