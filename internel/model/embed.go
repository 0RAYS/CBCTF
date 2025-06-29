package model

import (
	"CBCTF/internel/config"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"net"
	"slices"
	"strings"
	"time"
)

type AvatarURL string

func (a AvatarURL) Value() (driver.Value, error) {
	if a == "" {
		return nil, nil
	}
	return strings.TrimPrefix(string(a), config.Env.Backend), nil
}

func (a *AvatarURL) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan AvatarURL: %v", value)
	}
	if len(bytes) == 0 {
		*a = ""
		return nil
	}
	*a = AvatarURL(config.Env.Backend + string(bytes))
	return nil
}

type StringList []string

func (s StringList) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringList value")
	}
	if len(bytes) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(bytes, s)
}

type UintList []uint

func (u UintList) Value() (driver.Value, error) {
	if len(u) == 0 {
		return nil, nil
	}
	return json.Marshal(u)
}

func (u *UintList) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan UintList value")
	}
	if len(bytes) == 0 {
		*u = nil
		return nil
	}
	return json.Unmarshal(bytes, u)
}

type StringMap map[string]string

func (s StringMap) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}
	return json.Marshal(s)
}

func (s *StringMap) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan StringMap value")
	}
	if len(bytes) == 0 {
		*s = nil
		return nil
	}
	return json.Unmarshal(bytes, s)
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

type Target struct {
	Hostname string   `json:"hostname"`
	CIDR     string   `json:"cidr"`
	Except   []string `json:"except"`
}

func (t Target) isValidIPBlock() bool {
	if t.Hostname != "" {
		return true
	}
	if t.Hostname == "" && t.CIDR == "" {
		return false
	}
	_, ipNet, err := net.ParseCIDR(t.CIDR)
	if err != nil {
		return false
	}
	for _, ex := range t.Except {
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
		p.From = slices.DeleteFunc(p.From, func(t Target) bool {
			return !t.isValidIPBlock()
		})
		p.To = slices.DeleteFunc(p.From, func(t Target) bool {
			return !t.isValidIPBlock()
		})
	}
	if len(n) == 0 {
		return nil, nil
	}
	return json.Marshal(n)
}

func (n *NetworkPolicies) Scan(value any) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan NetworkPolicy value")
	}
	if len(bytes) == 0 {
		*n = nil
		return nil
	}
	return json.Unmarshal(bytes, n)
}

type Ports []int32

func (e Ports) Value() (driver.Value, error) {
	tmp := make([]int32, 0)
	for _, port := range e {
		if port > 0 && port < 65535 {
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
