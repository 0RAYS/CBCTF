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

func (s *Strings) Scan(value any) error {
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

func (u *Uints) Scan(value any) error {
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

func (p *Prizes) Scan(value any) error {
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

func (t *Timelines) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Timelines value")
	}
	return json.Unmarshal(bytes, t)
}

type Target struct {
	Hostname string   `json:"hostname"`
	CIDR     string   `json:"cidr"`
	Except   []string `json:"except"`
}

func isValidIPBlock(ipBlock Target) bool {
	if ipBlock.Hostname != "" {
		return true
	}
	if ipBlock.Hostname == "" && ipBlock.CIDR == "" {
		return false
	}
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
	From []Target `json:"from"`
	To   []Target `json:"to"`
}

var DefaultNetworkPolicy = NetworkPolicy{
	From: []Target{},
	To: []Target{
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

func (n *NetworkPolicies) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	return json.Unmarshal(bytes, n)
}

// Docker 题目的 Docker 配置, 一个容器可以有多个 flag 和多个映射端口
type Docker struct {
	PodGroup        uint            `json:"pod_group"`
	Hostname        string          `json:"hostname"`
	FlagIDL         []uint          `json:"flag_id"`
	Flags           Strings         `json:"flags"`
	Image           string          `json:"image"`
	Ports           Ports           `json:"ports"`
	NetworkPolicies NetworkPolicies `json:"network_policies"`
}

type Dockers []Docker

func (d Dockers) Value() (driver.Value, error) {
	for i, docker := range d {
		for j, port := range docker.Ports {
			if port < 1 || port > 65535 {
				d[i].Ports = append(d[i].Ports[:j], d[i].Ports[j+1:]...)
			}
		}
		for j, ipBlock := range docker.NetworkPolicies {
			for k, from := range ipBlock.From {
				if !isValidIPBlock(from) {
					d[i].NetworkPolicies[j].From = append(d[i].NetworkPolicies[j].From[:k], d[i].NetworkPolicies[j].From[k+1:]...)
				}
			}
			for k, to := range ipBlock.To {
				if !isValidIPBlock(to) {
					d[i].NetworkPolicies[j].To = append(d[i].NetworkPolicies[j].To[:k], d[i].NetworkPolicies[j].To[k+1:]...)
				}
			}
		}
	}
	return json.Marshal(d)
}

func (d *Dockers) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Dockers value")
	}
	return json.Unmarshal(bytes, d)
}

type Ports []int32

func (e Ports) Value() (driver.Value, error) {
	tmp := make([]int32, 0)
	for _, port := range e {
		if port > 1 && port < 65535 {
			tmp = append(tmp, port)
		}
	}
	return json.Marshal(e)
}

func (e *Ports) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Ports value")
	}
	return json.Unmarshal(bytes, e)
}

type References struct {
	UserID    uint `json:"user_id"`
	TeamID    uint `json:"team_id"`
	ContestID uint `json:"contest_id"`
	UsageID   uint `json:"usage_id"`
}

func (r References) Value() (driver.Value, error) {
	return json.Marshal(r)
}

func (r *References) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan References value")
	}
	return json.Unmarshal(bytes, r)
}
